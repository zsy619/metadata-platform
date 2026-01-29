package metadata

import (
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/assert"
)

// TestRouterRegistration 验证路由注册逻辑是否在代码层生效
// 执行命令: go test -v apps/backend/internal/module/metadata/router_test.go
func TestRouterRegistration(t *testing.T) {
	h := server.Default()

	// 模拟所需的参数 (即使为 nil，只要 RegisterRoutes 不在注册时立即调用它们即可)
	// 由于我们的 RegisterRoutes 主要是闭包注册，所以只要不触发 Handler 逻辑，注入 nil 是安全的
	RegisterRoutes(h, nil, nil, nil)

	// 1. 验证探测接口 /api/metadata/ping
	w1 := ut.PerformRequest(h.Engine, "GET", "/api/metadata/ping", nil)
	resp1 := w1.Result()
	assert.Equal(t, http.StatusOK, resp1.StatusCode())
	assert.Contains(t, string(resp1.Body()), "pong")

	// 2. 验证工具接口 (16, 32, 64 位)
	testRoutes := []string{
		"/api/metadata/utils/generate-model-code-16",
		"/api/metadata/utils/generate-model-code",
		"/api/metadata/utils/generate-model-code-64",
	}

	for _, route := range testRoutes {
		w := ut.PerformRequest(h.Engine, "GET", route, nil)
		resp := w.Result()
		// 只要不返回 404，就证明路由匹配已成功
		assert.NotEqual(t, http.StatusNotFound, resp.StatusCode(), "Route %s should be registered", route)
	}
}
