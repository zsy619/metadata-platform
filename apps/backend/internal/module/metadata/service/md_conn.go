package service

import (
	"errors"

	"fmt"

	"metadata-platform/internal/module/metadata/adapter"
	"metadata-platform/internal/module/metadata/model"
	"metadata-platform/internal/module/metadata/repository"
	"metadata-platform/internal/utils"
)

// MdConnService 数据连接服务接口
type MdConnService interface {
	CreateConn(conn *model.MdConn) error
	GetConnByID(id string) (*model.MdConn, error)
	GetConnByName(name string) (*model.MdConn, error)
	UpdateConn(conn *model.MdConn) error
	DeleteConn(id string) error
	GetAllConns(tenantID string) ([]model.MdConn, error)
	GetConnsByParentID(parentID string) ([]model.MdConn, error)
	TestConnection(conn *model.MdConn) error
	GetTables(conn *model.MdConn, schema string) ([]adapter.TableInfo, error)
	GetViews(conn *model.MdConn, schema string) ([]adapter.ViewInfo, error)
	GetTableStructure(conn *model.MdConn, schema, table string) ([]adapter.ColumnInfo, error)
	PreviewTableData(conn *model.MdConn, schema, table string, limit int) ([]map[string]interface{}, error)
	GetSchemas(conn *model.MdConn) ([]string, error)
	ExecuteSQLForColumns(conn *model.MdConn, query string, params map[string]interface{}) ([]adapter.ColumnInfo, error)
}

// mdConnService 数据连接服务实现
type mdConnService struct {
	connRepo   repository.MdConnRepository
	snowflake  *utils.Snowflake
}

// NewMdConnService 创建数据连接服务实例
func NewMdConnService(connRepo repository.MdConnRepository) MdConnService {
	// 创建雪花算法生成器实例，使用默认数据中心ID和机器ID
	snowflake := utils.NewSnowflake(1, 1)
	return &mdConnService{
		connRepo:  connRepo,
		snowflake: snowflake,
	}
}

// CreateConn 创建数据连接
func (s *mdConnService) CreateConn(conn *model.MdConn) error {
	// 检查数据连接名称是否已存在
	existingConn, err := s.connRepo.GetConnByName(conn.ConnName)
	if err == nil && existingConn != nil {
		return errors.New("数据连接名称已存在")
	}

	// 使用雪花算法生成唯一ID
	conn.ID = s.snowflake.GenerateIDString()

	// 创建数据连接
	return s.connRepo.CreateConn(conn)
}

// GetConnByID 根据ID获取数据连接
func (s *mdConnService) GetConnByID(id string) (*model.MdConn, error) {
	return s.connRepo.GetConnByID(id)
}

// GetConnByName 根据名称获取数据连接
func (s *mdConnService) GetConnByName(name string) (*model.MdConn, error) {
	return s.connRepo.GetConnByName(name)
}

// UpdateConn 更新数据连接
func (s *mdConnService) UpdateConn(conn *model.MdConn) error {
	// 检查数据连接是否存在
	existingConn, err := s.connRepo.GetConnByID(conn.ID)
	if err != nil {
		return errors.New("数据连接不存在")
	}

	// 如果数据连接名称发生变化，检查新名称是否已存在
	if existingConn.ConnName != conn.ConnName {
		otherConn, err := s.connRepo.GetConnByName(conn.ConnName)
		if err == nil && otherConn != nil {
			return errors.New("数据连接名称已存在")
		}
	}

	// 更新数据连接
	return s.connRepo.UpdateConn(conn)
}

// DeleteConn 删除数据连接
func (s *mdConnService) DeleteConn(id string) error {
	// 检查数据连接是否存在
	_, err := s.connRepo.GetConnByID(id)
	if err != nil {
		return errors.New("数据连接不存在")
	}

	// 删除数据连接
	return s.connRepo.DeleteConn(id)
}

// GetAllConns 获取所有数据连接
func (s *mdConnService) GetAllConns(tenantID string) ([]model.MdConn, error) {
	return s.connRepo.GetAllConns(tenantID)
}

// GetConnsByParentID 根据父ID获取数据连接
func (s *mdConnService) GetConnsByParentID(parentID string) ([]model.MdConn, error) {
	return s.connRepo.GetConnsByParentID(parentID)
}

// getExtractor 根据数据连接获取元数据提取器
func (s *mdConnService) getExtractor(conn *model.MdConn) (adapter.MetadataExtractor, error) {
	switch conn.ConnKind {
	case "MySQL", "TiDB", "OceanBase", "Doris", "StarRocks":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnHost,
			conn.ConnPort,
			conn.ConnDatabase,
		)
		return adapter.NewMySQLExtractor(dsn)
	case "PostgreSQL", "OpenGauss", "Kingbase":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			conn.ConnHost,
			conn.ConnPort,
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnDatabase,
		)
		return adapter.NewPostgreSQLExtractor(dsn)
	case "SQL Server":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnHost,
			conn.ConnPort,
			conn.ConnDatabase,
		)
		return adapter.NewSQLServerExtractor(dsn)
	case "Oracle":
		dsn := fmt.Sprintf("%s/%s@%s:%d/%s",
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnHost,
			conn.ConnPort,
			conn.ConnDatabase,
		)
		return adapter.NewOracleExtractor(dsn)
	case "SQLite":
		// SQLite 使用文件路径作为 DSN
		dsn := conn.ConnHost // 假设 ConnHost 存储文件路径
		return adapter.NewSQLiteExtractor(dsn)
	case "ClickHouse":
		dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s",
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnHost,
			conn.ConnPort,
			conn.ConnDatabase,
		)
		return adapter.NewClickHouseExtractor(dsn)
	case "DM":
		dsn := fmt.Sprintf("dm://%s:%s@%s:%d",
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnHost,
			conn.ConnPort,
		)
		return adapter.NewDamengExtractor(dsn)
	case "MongoDB":
		dsn := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnHost,
			conn.ConnPort,
			conn.ConnDatabase,
		)
		return adapter.NewMongoDBExtractor(dsn)
	case "Redis":
		dsn := fmt.Sprintf("redis://%s:%s@%s:%d/%s",
			conn.ConnUser,
			conn.ConnPassword,
			conn.ConnHost,
			conn.ConnPort,
			conn.ConnDatabase,
		)
		return adapter.NewRedisExtractor(dsn)
	default:
		return nil, errors.New("不支持的数据源类型: " + conn.ConnKind)
	}
}

// TestConnection 测试数据连接
func (s *mdConnService) TestConnection(conn *model.MdConn) error {
	extractor, err := s.getExtractor(conn)
	if err != nil {
		// 获取提取器失败，更新状态为 2 (连接失败)
		conn.State = 2
		s.connRepo.UpdateConn(conn)
		return err
	}
	defer extractor.Close()

	// 测试连接
	if err := extractor.TestConnection(); err != nil {
		// 测试失败，更新状态为 2 (连接失败)
		conn.State = 2
		if updateErr := s.connRepo.UpdateConn(conn); updateErr != nil {
			fmt.Printf("Warning: Failed to update connection state to error: %v\n", updateErr)
		}
		return err
	}

	// 测试成功后，更新数据库中的状态为 1 (有效)
	conn.State = 1
	if updateErr := s.connRepo.UpdateConn(conn); updateErr != nil {
		// 记录错误但不影响测试结果
		fmt.Printf("Warning: Failed to update connection state: %v\n", updateErr)
	}

	return nil
}

// GetTables 获取数据库表列表
func (s *mdConnService) GetTables(conn *model.MdConn, schema string) ([]adapter.TableInfo, error) {
	extractor, err := s.getExtractor(conn)
	if err != nil {
		return nil, err
	}
	defer extractor.Close()

	if schema == "" {
		schema = conn.ConnDatabase
	}
	return extractor.GetTables(schema)
}

// GetViews 获取数据库视图列表
func (s *mdConnService) GetViews(conn *model.MdConn, schema string) ([]adapter.ViewInfo, error) {
	extractor, err := s.getExtractor(conn)
	if err != nil {
		return nil, err
	}
	defer extractor.Close()

	if schema == "" {
		schema = conn.ConnDatabase
	}
	return extractor.GetViews(schema)
}

// GetTableStructure 获取表结构
func (s *mdConnService) GetTableStructure(conn *model.MdConn, schema, table string) ([]adapter.ColumnInfo, error) {
	extractor, err := s.getExtractor(conn)
	if err != nil {
		return nil, err
	}
	defer extractor.Close()

	if schema == "" {
		schema = conn.ConnDatabase
	}
	return extractor.GetColumns(schema, table)
}

// PreviewTableData 预览表数据
func (s *mdConnService) PreviewTableData(conn *model.MdConn, schema, table string, limit int) ([]map[string]interface{}, error) {
	extractor, err := s.getExtractor(conn)
	if err != nil {
		return nil, err
	}
	defer extractor.Close()

	if schema == "" {
		schema = conn.ConnDatabase
	}
	return extractor.PreviewData(schema, table, limit)
}

// GetSchemas 获取数据库模式列表
func (s *mdConnService) GetSchemas(conn *model.MdConn) ([]string, error) {
	extractor, err := s.getExtractor(conn)
	if err != nil {
		return nil, err
	}
	defer extractor.Close()
	return extractor.GetSchemas()
}

// ExecuteSQLForColumns 执行SQL并获取返回列信息
func (s *mdConnService) ExecuteSQLForColumns(conn *model.MdConn, query string, params map[string]interface{}) ([]adapter.ColumnInfo, error) {
	extractor, err := s.getExtractor(conn)
	if err != nil {
		return nil, err
	}
	defer extractor.Close()

	// 简单的参数处理：将 map 转换为 list (如果有顺序要求，这里需要更复杂的解析，目前假设无参数或顺序参数)
	// 实际应用中，如果 SQL 使用 ? 或 $1 占位符，params 需要按顺序提供
	// 如果是命名参数 @name，需要 SQL 驱动支持或手动替换
	// 这里暂且只支持无参或简单按序 (需前端保证)
	var args []interface{}
	// TODO: 实现命名参数解析
	
	return extractor.GetQueryColumns(query, args)
}
