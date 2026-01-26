package api

import (
	"context"
	"fmt"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DynamicRouter 动态路由注册器
type DynamicRouter struct {
	hertz   *server.Hertz
	svc     *service.Services
}

// NewDynamicRouter 创建动态路由注册器实例
func NewDynamicRouter(hertz *server.Hertz, svc *service.Services) *DynamicRouter {
	return &DynamicRouter{
		hertz: hertz,
		svc:   svc,
	}
}

// LoadAndRegisterAll 加载并注册所有活跃的动态接口
func (r *DynamicRouter) LoadAndRegisterAll() error {
	apis, err := r.svc.API.GetAllAPIs()
	if err != nil {
		return err
	}

	for _, a := range apis {
		if a.State == 1 { // 仅注册启用状态的接口
			r.RegisterAPI(a.Method, a.Path, a.Code)
		}
	}

	return nil
}

// RegisterAPI 注册单个动态路由
func (r *DynamicRouter) RegisterAPI(method, path, apiCode string) {
	utils.SugarLogger.Infof("Registering dynamic route: [%s] %s (code: %s)", method, path, apiCode)
	
	// 绑定通用处理器
	// 注意：apiCode 在这里被用作查找 modelID 的依据（目前简单约定 apiCode = modelCode + "_" + method）
	// 后续可以增加 MdModelAPI 关联表来更精确对应
	
	handler := r.getGenericHandler(apiCode)
	r.hertz.Handle(method, path, handler)
}

func (r *DynamicRouter) getGenericHandler(apiCode string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 1. 获取 ModelID (解析 apiCode)
		modelCode := ""
		if idx := lastIndex(apiCode, "_"); idx != -1 {
			modelCode = apiCode[:idx]
		} else {
			modelCode = apiCode
		}

		md, err := r.svc.Model.GetModelByCode(modelCode)
		if err != nil || md == nil {
			utils.ErrorResponse(ctx, consts.StatusNotFound, "模型配置解析失败")
			return
		}

		method := string(ctx.Method())

		// 2. 根据方法分发逻辑
		switch method {
		case "POST":
			var data map[string]any
			if err := ctx.Bind(&data); err != nil {
				utils.ErrorResponse(ctx, consts.StatusBadRequest, "参数解析失败")
				return
			}
			res, err := r.svc.CRUD.Create(md.ID, data)
			if err != nil {
				utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
				return
			}
			utils.SuccessResponse(ctx, res)

		case "GET":
			id := ctx.Param("id")
			if id != "" {
				res, err := r.svc.CRUD.Get(md.ID, id)
				if err != nil {
					utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
					return
				}
				if res == nil {
					utils.ErrorResponse(ctx, consts.StatusNotFound, "记录不存在")
					return
				}
				utils.SuccessResponse(ctx, res)
			} else {
				res, count, err := r.svc.CRUD.List(md.ID, nil)
				if err != nil {
					utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
					return
				}
				utils.SuccessResponse(ctx, map[string]any{
					"list":  res,
					"total": count,
				})
			}

		case "PUT":
			id := ctx.Param("id")
			var data map[string]any
			if err := ctx.Bind(&data); err != nil {
				utils.ErrorResponse(ctx, consts.StatusBadRequest, "参数解析失败")
				return
			}
			err := r.svc.CRUD.Update(md.ID, id, data)
			if err != nil {
				utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
				return
			}
			utils.SuccessResponse(ctx, nil)

		case "DELETE":
			id := ctx.Param("id")
			err := r.svc.CRUD.Delete(md.ID, id)
			if err != nil {
				utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
				return
			}
			utils.SuccessResponse(ctx, nil)

		default:
			utils.ErrorResponse(ctx, consts.StatusMethodNotAllowed, fmt.Sprintf("Unsupported method: %s", method))
		}
	}
}

func lastIndex(s, sep string) int {
	for i := len(s) - len(sep); i >= 0; i-- {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}
