package service

import (
	"context"
	"errors"
	"fmt"

	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
)

// MasterDetailService 主子表服务接口
type MasterDetailService interface {
	CreateMasterDetail(ctx context.Context, masterModelID string, detailModelID string, payload map[string]any) error
}

type masterDetailService struct {
	crudSvc      CRUDService
	relationRepo repository.MdModelRelationRepository
	modelRepo    repository.MdModelRepository
	executor     *engine.SQLExecutor
}

// NewMasterDetailService 创建主子表服务实例
func NewMasterDetailService(
	crudSvc CRUDService,
	relationRepo repository.MdModelRelationRepository,
	modelRepo repository.MdModelRepository,
	executor *engine.SQLExecutor,
) MasterDetailService {
	return &masterDetailService{
		crudSvc:      crudSvc,
		relationRepo: relationRepo,
		modelRepo:    modelRepo,
		executor:     executor,
	}
}

// CreateMasterDetail 创建主子表数据
func (s *masterDetailService) CreateMasterDetail(ctx context.Context, masterModelID string, detailModelID string, payload map[string]any) error {
	// 1. 获取模型和关系定义
	relation, err := s.relationRepo.GetRelation(masterModelID, detailModelID)
	if err != nil {
		return fmt.Errorf("failed to get relation: %v", err)
	}
	if relation == nil {
		return fmt.Errorf("relation between %s and %s not found", masterModelID, detailModelID)
	}

	masterData, ok := payload["master"].(map[string]any)
	if !ok {
		return errors.New("invalid master data format")
	}

	var detailList []map[string]any
	if dList, ok := payload["details"].([]any); ok {
		for _, item := range dList {
			if d, ok := item.(map[string]any); ok {
				detailList = append(detailList, d)
			}
		}
	} else if dList, ok := payload["details"].([]map[string]any); ok {
		detailList = dList
	}

	if len(detailList) == 0 {
		// 允许只创建主表? 暂时允许
	}

	// 2. 获取主表数据库连接信息
	masterModel, err := s.modelRepo.GetModelByID(masterModelID)
	if err != nil {
		return err
	}

	db, err := s.executor.GetConnection(masterModel.ConnID)
	if err != nil {
		return err
	}

	// 3. 开启事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return tx.Error
	}

	// 4. 创建主表数据
	createdMaster, err := s.crudSvc.CreateWithTx(ctx, masterModelID, masterData, tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create master record: %v", err)
	}

	// 5. 获取主键值
	// 假设 dataValidator 或者 builder 能够识别主键，但是 Create 返回的是 data map。
	// 如果是自增ID，CreateWithTx 目前的实现可能没有把 ID 回填到 map 中（如果是 executor 执行的 insert）。
	// 目前的 SQLBuilder 生成的 insert 语句 和 SQLExecutor 的 execution 并没有自动获取 LastInsertId 并回填。

	// 这里是一个潜在的坑: 如果主键是自增的，我们需要获取它。
	// 如果主键是 UUID (客户端生成或者 validator 生成)，那么 map 里已经有了。

	// 为了演示完整性，我们假设主键在 masterData 中已经存在（例如 UUID），或者底层支持回填。
	// 如果需要支持自增，Metadata模块应该配置为 "服务端生成UUID" 或 "客户端提供"。
	// 暂时假设: 主键字段名为 'id' 且已在 createdMaster 中 (如果 Create 只有 data copy, 那么必须在 data 中)

	// 在实际生产中，Executor 应该 Scan 插入后的 ID，或者使用 RETURNING (Postgres) / LastInsertId (MySQL)。
	// 目前的 Execute 方法没有返回 ID。
	// 这是一个已知限制，我们在 Phase 3 应该考虑增强 SQLExecutor 或者是要求 UUID。

	var masterID string
	if id, ok := createdMaster["id"]; ok {
		masterID = utils.ToString(id)
	} else {
		// 尝试查找任意 ID 字段?
		// 暂时无法获取自增ID，回滚
		tx.Rollback()
		return errors.New("failed to retrieve master ID (auto-increment not supported in current master-detail impl, please use UUID)")
	}

	// 6. 创建子表数据
	if len(detailList) > 0 {
		for i := range detailList {
			// 设置外键
			detailList[i][relation.ForeignKey] = masterID
		}

		_, err := s.crudSvc.BatchCreateWithTx(ctx, detailModelID, detailList, tx)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create detail records: %v", err)
		}
	}

	// 7. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	return nil
}
