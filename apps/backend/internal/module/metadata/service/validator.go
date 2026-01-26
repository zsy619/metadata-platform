package service

import (
	"fmt"
	"metadata-platform/internal/module/metadata/model"
	"regexp"
	"strconv"
)

// DataValidator 数据校验器接口
type DataValidator interface {
	Validate(modelID string, fields []*model.MdModelField, data map[string]any) error
}

type dataValidator struct{}

// NewDataValidator 创建数据校验器实例
func NewDataValidator() DataValidator {
	return &dataValidator{}
}

// Validate 执行数据校验
func (v *dataValidator) Validate(modelID string, fields []*model.MdModelField, data map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	// 建立字段映射
	fieldMap := make(map[string]*model.MdModelField)
	for _, f := range fields {
		fieldMap[f.ColumnName] = f
	}

	// 1. 检查是否存在未定义字段（只在此处报错，如果需要严格模式）
	// 目前策略：忽略未定义字段，或者根据配置报错
	for k := range data {
		if _, ok := fieldMap[k]; !ok {
			// 如果是主键 id 且正在做 Create，可能由系统生成，目前暂允许
			if k != "id" {
				return fmt.Errorf("字段 '%s' 未在模型 %s 中定义", k, modelID)
			}
		}
	}

	// 2. 遍历元数据定义，执行逐项校验
	for _, f := range fields {
		val, exists := data[f.ColumnName]

		// 2.1 必填项校验 (非空且非自增)
		if !f.IsNullable && !f.IsAutoIncrement {
			if !exists || val == nil || fmt.Sprintf("%v", val) == "" {
				return fmt.Errorf("字段 '%s' (%s) 不能为空", f.ColumnName, f.ShowTitle)
			}
		}

		// 如果数据不存在或为空，后续校验跳过
		if !exists || val == nil || fmt.Sprintf("%v", val) == "" {
			continue
		}

		// 2.2 类型与长度/范围校验
		strVal := fmt.Sprintf("%v", val)
		
		switch f.FieldType {
		case "string":
			if f.MaxLength > 0 && len(strVal) > f.MaxLength {
				return fmt.Errorf("字段 '%s' 长度不能超过 %d (当前: %d)", f.ColumnName, f.MaxLength, len(strVal))
			}
		case "integer", "long", "decimal":
			num, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				return fmt.Errorf("字段 '%s' 格式不正确，应为数值", f.ColumnName)
			}
			if f.Max > f.Min { // 仅当设置了有效区间时校验
				if num > f.Max {
					return fmt.Errorf("字段 '%s' 数值不能大于 %v", f.ColumnName, f.Max)
				}
				if num < f.Min {
					return fmt.Errorf("字段 '%s' 数值不能小于 %v", f.ColumnName, f.Min)
				}
			}
		}

		// 2.3 正则校验
		if f.ValidationRule != "" {
			matched, err := regexp.MatchString(f.ValidationRule, strVal)
			if err != nil {
				return fmt.Errorf("字段 '%s' 的校验规则 (正则) 配置错误: %v", f.ColumnName, err)
			}
			if !matched {
				return fmt.Errorf("字段 '%s' 不符合业务校验规则", f.ColumnName)
			}
		}
	}

	return nil
}
