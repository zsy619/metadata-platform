package utils

import "fmt"

// ToString 将 any 类型安全转换为 string
func ToString(v any) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	default:
		return fmt.Sprintf("%v", v)
	}
}
