package api

import (
	"context"
	"fmt"
	"metadata-platform/internal/module/metadata/api/middleware"
	"metadata-platform/internal/module/metadata/service"
	"metadata-platform/internal/utils"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// DynamicRouter 动态路由注册器
type DynamicRouter struct {
	*utils.BaseHandler
	hertz        *server.Hertz
	svc          *service.Services
	queryHandler *DataQueryHandler
}

// NewDynamicRouter 创建动态路由注册器实例
func NewDynamicRouter(hertz *server.Hertz, svc *service.Services) *DynamicRouter {
	return &DynamicRouter{
		BaseHandler:  utils.NewBaseHandler(),
		hertz:        hertz,
		svc:          svc,
		queryHandler: NewDataQueryHandler(svc.CRUD, svc.Model),
	}
}

// LoadAndRegisterAll 加载并注册所有活跃的动态接口
func (r *DynamicRouter) LoadAndRegisterAll() error {
	apis, err := r.svc.API.GetAllAPIs()
	if err != nil {
		return err
	}

	for _, a := range apis {
		if a.Status == 1 { // 仅注册启用状态的接口
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

	// 添加审计中间件
	auditMiddleware := middleware.AuditMiddleware(r.svc.Audit)

	r.hertz.Handle(method, path, auditMiddleware, handler)
}

func (r *DynamicRouter) getGenericHandler(apiCode string) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 1. 获取 ModelID (解析 apiCode)
		var modelCode string

		// 优先匹配特殊后缀
		var handlerType string
		if before, ok := strings.CutSuffix(apiCode, "_QUERY"); ok {
			modelCode = before
			handlerType = "QUERY"
		} else if before, ok := strings.CutSuffix(apiCode, "_BATCH_CREATE"); ok {
			modelCode = before
			handlerType = "BATCH_CREATE"
		} else if before, ok := strings.CutSuffix(apiCode, "_BATCH_DELETE"); ok {
			modelCode = before
			handlerType = "BATCH_DELETE"
		} else if before, ok := strings.CutSuffix(apiCode, "_STATISTICS"); ok {
			modelCode = before
			handlerType = "STATISTICS"
		} else if before, ok := strings.CutSuffix(apiCode, "_AGGREGATE"); ok {
			modelCode = before
			handlerType = "AGGREGATE"
		} else {
			// 默认逻辑：取最后一个下划线前缀
			if idx := lastIndex(apiCode, "_"); idx != -1 {
				modelCode = apiCode[:idx]
			} else {
				modelCode = apiCode
			}
		}

		md, err := r.svc.Model.GetModelByCode(modelCode)
		if err != nil || md == nil {
			utils.ErrorResponse(ctx, consts.StatusNotFound, "模型配置解析失败")
			return
		}

		method := string(ctx.Method())

		// 特殊处理器分发
		switch handlerType {
		case "QUERY":
			r.queryHandler.HandleUnifiedQueryWithModelID(ctx, md.ID)
			return
		case "BATCH_CREATE":
			r.queryHandler.HandleBatchCreateWithModelID(ctx, md.ID)
			return
		case "BATCH_DELETE":
			r.queryHandler.HandleBatchDeleteWithModelID(ctx, md.ID)
			return
		case "STATISTICS":
			r.queryHandler.HandleStatisticsWithModelID(ctx, md.ID)
			return
		case "AGGREGATE":
			r.queryHandler.HandleAggregateWithModelID(ctx, md.ID)
			return
		}

		// 2. 根据方法分发逻辑
		switch method {
		case "POST":
			var data map[string]any
			if err := ctx.Bind(&data); err != nil {
				utils.ErrorResponse(ctx, consts.StatusBadRequest, "参数解析失败")
				return
			}
			// 从 Handler 中获取 Context
			// Hertz HandlerFunc: func(c context.Context, ctx *app.RequestContext)
			// 但是 DynamicRouter 的 HandleCreate 签名目前是 func(c context.Context, ctx *app.RequestContext)
			// 所以直接使用 c 即可

			result, err := r.svc.CRUD.Create(c, md.ID, data)
			if err != nil {
				utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
				return
			}
			utils.SuccessResponse(ctx, result)

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
			if err := r.svc.CRUD.Update(c, md.ID, id, data); err != nil {
				utils.ErrorResponse(ctx, consts.StatusInternalServerError, err.Error())
				return
			}
			utils.SuccessResponse(ctx, nil)

		case "DELETE":
			id := ctx.Param("id")
			if err := r.svc.CRUD.Delete(c, md.ID, id); err != nil {
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
