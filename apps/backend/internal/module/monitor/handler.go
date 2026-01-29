package monitor

import (
	"context"
	"math/rand"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/websocket"
)

var upgrader = websocket.HertzUpgrader{
	CheckOrigin: func(ctx *app.RequestContext) bool {
		return true
	},
}

// RegisterRoutes 注册监控模块路由
func RegisterRoutes(h *server.Hertz) {
	group := h.Group("/api/monitor")
	group.GET("/ws", WSHandler)
}

// WSHandler WebSocket处理函数
func WSHandler(ctx context.Context, c *app.RequestContext) {
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		defer conn.Close()
		for {
			// 模拟实时数据
			data := map[string]interface{}{
				"timestamp": time.Now().Format("15:04:05"),
				"requests":  1000 + rand.Intn(500),
				"qps":       300 + rand.Intn(200),
				"error_rate": float64(rand.Intn(10)) / 100.0,
				"latency":    50 + rand.Intn(100),
				"type": "realtime",
			}
			
			if err := conn.WriteJSON(data); err != nil {
				break
			}
			time.Sleep(2 * time.Second)
		}
	})
	if err != nil {
		return
	}
}
