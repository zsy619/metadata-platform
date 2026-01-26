package service

import (
	"context"
	"errors"
	"fmt"
	"metadata-platform/internal/module/metadata/engine"
	"metadata-platform/internal/module/metadata/repository"
)

// TreeService 树形结构服务接口
type TreeService interface {
	GetTree(modelID string) ([]map[string]any, error)
	GetChildren(modelID string, parentID string) ([]map[string]any, error)
	GetPath(modelID string, id string) ([]map[string]any, error)
	AddNode(modelID string, data map[string]any) (map[string]any, error)
	MoveNode(modelID string, id string, targetParentID string) error
	DeleteNode(modelID string, id string) error
}

type treeService struct {
	modelRepo repository.MdModelRepository
	crudSvc   CRUDService
	executor  *engine.SQLExecutor
}

// NewTreeService 创建树形结构服务实例
func NewTreeService(modelRepo repository.MdModelRepository, crudSvc CRUDService, executor *engine.SQLExecutor) TreeService {
	return &treeService{
		modelRepo: modelRepo,
		crudSvc:   crudSvc,
		executor:  executor,
	}
}

// GetTree 获取完整树结构
func (s *treeService) GetTree(modelID string) ([]map[string]any, error) {
	md, err := s.modelRepo.GetModelByID(modelID)
	if err != nil {
		return nil, err
	}
	if !md.IsTree {
		return nil, errors.New("model is not configured as tree structure")
	}

	// 获取所有数据
	// TODO: 大数据量时应考虑懒加载，这里假设数据量适中，一次性加载
	list, _, err := s.crudSvc.List(modelID, map[string]any{
		"page_size": 10000,
	})
	if err != nil {
		return nil, err
	}

	return s.buildTree(list, "0", md.TreeParentField, "id")
}

// GetChildren 获取直接子节点
func (s *treeService) GetChildren(modelID string, parentID string) ([]map[string]any, error) {
	md, err := s.modelRepo.GetModelByID(modelID)
	if err != nil {
		return nil, err
	}
	if !md.IsTree {
		return nil, errors.New("model is not configured as tree structure")
	}

	// 构造查询参数
	queryParams := map[string]any{
		"filters": []any{
			map[string]any{
				"column_name": md.TreeParentField,
				"operator":    "=",
				"value":       parentID,
			},
		},
	}

	list, _, err := s.crudSvc.List(modelID, queryParams)
	return list, err
}

// GetPath 获取节点路径
func (s *treeService) GetPath(modelID string, id string) ([]map[string]any, error) {
	// 简单的实现：向上递归查找
	// 优化实现：如果表中存了 path 字段，直接解析 path
	md, err := s.modelRepo.GetModelByID(modelID)
	if err != nil {
		return nil, err
	}

	path := make([]map[string]any, 0)
	currentID := id

	for currentID != "0" && currentID != "" {
		node, err := s.crudSvc.Get(modelID, currentID)
		if err != nil || node == nil {
			break
		}
		// prepend
		path = append([]map[string]any{node}, path...)

		if pID, ok := node[md.TreeParentField]; ok {
			currentID = fmt.Sprintf("%v", pID)
		} else {
			break
		}
	}

	return path, nil
}

// AddNode 添加节点
func (s *treeService) AddNode(modelID string, data map[string]any) (map[string]any, error) {
	md, err := s.modelRepo.GetModelByID(modelID)
	if err != nil {
		return nil, err
	}
	if !md.IsTree {
		// 非树形模型，直接调用普通创建
		return s.crudSvc.Create(context.Background(), modelID, data)
	}

	// 1. 自动计算 path 和 level (如果配置了字段)
	parentID := "0"
	if p, ok := data[md.TreeParentField]; ok {
		parentID = fmt.Sprintf("%v", p)
	}

	if md.TreePathField != "" || md.TreeLevelField != "" {
		// parentPath := "/"
		// parentLevel := 0

		if parentID != "0" {
			parent, err := s.crudSvc.Get(modelID, parentID)
			if err != nil {
				return nil, err
			}
			if parent != nil {
				// if md.TreePathField != "" {
				// 	if pp, ok := parent[md.TreePathField].(string); ok {
				// 		parentPath = pp
				// 	}
				// }
				// if md.TreeLevelField != "" {
				// 	// 假设 level 是 int
				// 	// 这里需要根据实际类型处理，简单处理先略过类型断言错误
				// 	parentLevel = 1 // default
				// }
			}
		}

		// 暂时无法在创建前知道 ID (除非使用 UUID 并且预生成)
		// 这里的 Path 计算策略需要调整：
		// 策略 A: 插入后更新 Path (需要两次 DB 操作)
		// 策略 B: Path 不包含自身 ID (Parent Path)

		// 采用策略 A
	}

	// 2. 创建节点
	newNode, err := s.crudSvc.Create(context.Background(), modelID, data)
	if err != nil {
		return nil, err
	}

	// 3. 后置处理 Path (如果 ID 已生成)
	// TODO: 更新 Path

	return newNode, nil
}

// MoveNode 移动节点
func (s *treeService) MoveNode(modelID string, id string, targetParentID string) error {
	md, err := s.modelRepo.GetModelByID(modelID)
	if err != nil {
		return err
	}

	// 1. 循环引用检测
	if id == targetParentID {
		return errors.New("cannot move node to itself")
	}
	// 检查 targetParentID 是否是 id 的后代
	if targetParentID != "0" {
		path, err := s.GetPath(modelID, targetParentID)
		if err != nil {
			return err
		}
		for _, node := range path {
			if fmt.Sprintf("%v", node["id"]) == id { // 假设 ID 字段名为 id
				return errors.New("cannot move node to its descendant")
			}
		}
	}

	// 2. 更新 ParentID
	err = s.crudSvc.Update(context.Background(), modelID, id, map[string]any{
		md.TreeParentField: targetParentID,
	})

	// 3. 级联更新子节点 Path (如果使用了 Path 字段)
	// TODO

	return err
}

// DeleteNode 删除节点
func (s *treeService) DeleteNode(modelID string, id string) error {
	// 级联删除子节点
	children, err := s.GetChildren(modelID, id)
	if err != nil {
		return err
	}
	for _, child := range children {
		childID := fmt.Sprintf("%v", child["id"])
		if err := s.DeleteNode(modelID, childID); err != nil {
			return err
		}
	}

	return s.crudSvc.Delete(context.Background(), modelID, id)
}

func (s *treeService) buildTree(list []map[string]any, parentID string, parentField string, idField string) ([]map[string]any, error) {
	tree := make([]map[string]any, 0)
	for _, item := range list {
		pID := fmt.Sprintf("%v", item[parentField])
		if pID == parentID || (parentID == "0" && (pID == "" || pID == "<nil>")) {
			id := fmt.Sprintf("%v", item[idField])
			children, _ := s.buildTree(list, id, parentField, idField)
			if len(children) > 0 {
				item["children"] = children
			}
			tree = append(tree, item)
		}
	}
	return tree, nil
}
