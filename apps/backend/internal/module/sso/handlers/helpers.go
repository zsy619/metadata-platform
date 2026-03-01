package handlers

import (
	"encoding/json"
	"fmt"
	"time"
)

// generateSessionID 生成会话 ID
func generateSessionID() string {
	return fmt.Sprintf("sid-%d", time.Now().UnixNano())
}

// marshalExtraData 将额外数据编码为 JSON 字符串
func marshalExtraData(data map[string]any) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// unmarshalExtraData 从 JSON 字符串解码额外数据
func unmarshalExtraData(data string) (map[string]any, error) {
	if data == "" {
		return make(map[string]any), nil
	}
	
	var result map[string]any
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getStringFromExtraData 从额外数据获取字符串值
func getStringFromExtraData(extraData string, key string) string {
	data, err := unmarshalExtraData(extraData)
	if err != nil {
		return ""
	}
	
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}
