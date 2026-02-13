package queue

import (
	"sync"

	"metadata-platform/internal/module/audit/model"
	"metadata-platform/internal/utils"

	"gorm.io/gorm"
)

// LogType 日志类型
type LogType int

const (
	LogTypeOperation LogType = iota
	LogTypeDataChange
	LogTypeLogin
	LogTypeAccess
)

// LogEntry 日志条目包装器
type LogEntry struct {
	Type       LogType
	Operation  *model.SysOperationLog
	DataChange *model.SysDataChangeLog
	Login      *model.SysLoginLog
	Access     *model.SysAccessLog
}

// AuditLogQueue 审计日志队列
type AuditLogQueue struct {
	db          *gorm.DB
	logChan     chan LogEntry
	workerCount int
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// NewAuditLogQueue 创建审计日志队列
func NewAuditLogQueue(db *gorm.DB, bufferSize int, workerCount int) *AuditLogQueue {
	return &AuditLogQueue{
		db:          db,
		logChan:     make(chan LogEntry, bufferSize),
		workerCount: workerCount,
		stopChan:    make(chan struct{}),
	}
}

// Start 启动消费者Worker
func (q *AuditLogQueue) Start() {
	for i := 0; i < q.workerCount; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
	utils.SugarLogger.Infof("Audit log queue started with %d workers", q.workerCount)
}

// Stop 停止队列（优雅退出）
func (q *AuditLogQueue) Stop() {
	close(q.stopChan)
	q.wg.Wait()
	close(q.logChan) // Replace with better graceful shutdown if needed, strictly we should drain.
	// Actually, closing stopChan signals workers to stop accepting NEW items, 
	// but we might want to process remaining items in buffer.
	// For simplicity: process remaining items then exit.
	utils.SugarLogger.Info("Audit log queue stopped")
}

func (q *AuditLogQueue) worker(id int) {
	defer q.wg.Done()
	for {
		select {
		case entry, ok := <-q.logChan:
			if !ok {
				return
			}
			q.processLog(entry)
		case <-q.stopChan:
			// Process remaining items in channel
			for {
				select {
				case entry := <-q.logChan:
					q.processLog(entry)
				default:
					return
				}
			}
		}
	}
}

func (q *AuditLogQueue) processLog(entry LogEntry) {
	var err error
	switch entry.Type {
	case LogTypeOperation:
		err = q.db.Create(entry.Operation).Error
	case LogTypeDataChange:
		err = q.db.Create(entry.DataChange).Error
	case LogTypeLogin:
		err = q.db.Create(entry.Login).Error
	case LogTypeAccess:
		err = q.db.Create(entry.Access).Error
	}

	if err != nil {
		utils.SugarLogger.Errorf("Failed to write audit log: %v", err)
	}
}

// PushOperation 推送操作日志
func (q *AuditLogQueue) PushOperation(log *model.SysOperationLog) {
	select {
	case q.logChan <- LogEntry{Type: LogTypeOperation, Operation: log}:
	default:
		utils.SugarLogger.Warn("Audit log queue full, dropping operation log")
	}
}

// PushDataChange 推送数据变更日志
func (q *AuditLogQueue) PushDataChange(log *model.SysDataChangeLog) {
	select {
	case q.logChan <- LogEntry{Type: LogTypeDataChange, DataChange: log}:
	default:
		utils.SugarLogger.Warn("Audit log queue full, dropping data change log")
	}
}

// PushLogin 推送登录日志
func (q *AuditLogQueue) PushLogin(log *model.SysLoginLog) {
	select {
	case q.logChan <- LogEntry{Type: LogTypeLogin, Login: log}:
	default:
		utils.SugarLogger.Warn("Audit log queue full, dropping login log")
	}
}

// PushAccess 推送访问日志
func (q *AuditLogQueue) PushAccess(log *model.SysAccessLog) {
	select {
	case q.logChan <- LogEntry{Type: LogTypeAccess, Access: log}:
	default:
		utils.SugarLogger.Warn("Audit log queue full, dropping access log")
	}
}
