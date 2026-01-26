package middleware

import (
	"context"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"gorm.io/gorm"
)

var enforcer *casbin.Enforcer

// InitCasbin 初始化 Casbin Enforcer 并集成 GORM 适配器
func InitCasbin(db *gorm.DB, modelPath string) error {
	// 使用 GORM 适配器，指向 casbin_rule 表
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}

	enforcer, err = casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return err
	}

	// 加载策略
	return enforcer.LoadPolicy()
}

// GetEnforcer 获取 Casbin Enforcer 实例
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// CasbinMiddleware Casbin 权限校验中间件
func CasbinMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		if enforcer == nil {
			ctx.Next(c)
			return
		}

		// 从上下文中获取用户信息
		userID, exists := ctx.Get("user_id")
		if !exists {
			ctx.JSON(consts.StatusUnauthorized, map[string]any{
				"code":    401,
				"message": "未识别的用户身份",
			})
			ctx.Abort()
			return
		}

		// 管理员直接放通
		isAdmin, _ := ctx.Get("is_admin")
		if isAdmin != nil && isAdmin.(bool) {
			ctx.Next(c)
			return
		}

		// 获取请求路径和方法
		obj := string(ctx.Request.URI().Path())
		act := string(ctx.Method())
		sub := userID.(string)

		// 权限校验
		ok, err := enforcer.Enforce(sub, obj, act)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, map[string]any{
				"code":    500,
				"message": "权限校验异常: " + err.Error(),
			})
			ctx.Abort()
			return
		}

		if !ok {
			ctx.JSON(consts.StatusForbidden, map[string]any{
				"code":    403,
				"message": "权限不足",
			})
			ctx.Abort()
			return
		}

		ctx.Next(c)
	}
}
