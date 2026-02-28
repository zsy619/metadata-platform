package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"metadata-platform/internal/module/sso/service"
)

type SsoHandler struct {
	services *service.Services
}

func NewSsoHandler(services *service.Services) *SsoHandler {
	return &SsoHandler{
		services: services,
	}
}

func (h *SsoHandler) RegisterRoutes(r *server.Hertz) {
	ssoGroup := r.Group("/api/sso")
	{
		configGroup := ssoGroup.Group("/config")
		{
			configGroup.GET("", h.ListProtocolConfigs)
			configGroup.GET("/:id", h.GetProtocolConfig)
			configGroup.POST("", h.CreateProtocolConfig)
			configGroup.PUT("/:id", h.UpdateProtocolConfig)
			configGroup.DELETE("/:id", h.DeleteProtocolConfig)
			configGroup.GET("/type/:type", h.GetConfigsByType)
		}

		clientGroup := ssoGroup.Group("/client")
		{
			clientGroup.GET("", h.ListClients)
			clientGroup.GET("/:id", h.GetClient)
			clientGroup.POST("", h.CreateClient)
			clientGroup.PUT("/:id", h.UpdateClient)
			clientGroup.DELETE("/:id", h.DeleteClient)
		}

		keyGroup := ssoGroup.Group("/key")
		{
			keyGroup.GET("", h.ListKeys)
			keyGroup.GET("/:id", h.GetKey)
			keyGroup.POST("", h.CreateKey)
			keyGroup.POST("/generate", h.GenerateKey)
			keyGroup.PUT("/:id", h.UpdateKey)
			keyGroup.DELETE("/:id", h.DeleteKey)
		}

		mappingGroup := ssoGroup.Group("/mapping")
		{
			mappingGroup.GET("", h.ListFieldMappings)
			mappingGroup.GET("/:id", h.GetFieldMapping)
			mappingGroup.POST("", h.CreateFieldMapping)
			mappingGroup.PUT("/:id", h.UpdateFieldMapping)
			mappingGroup.DELETE("/:id", h.DeleteFieldMapping)
		}

		sessionGroup := ssoGroup.Group("/session")
		{
			sessionGroup.GET("", h.ListSessions)
			sessionGroup.GET("/:id", h.GetSession)
			sessionGroup.DELETE("/:id", h.DeleteSession)
			sessionGroup.POST("/:id/revoke", h.RevokeSession)
		}

		authGroup := ssoGroup.Group("/auth")
		{
			authGroup.GET("/:protocol/authorize", h.Authorize)
			authGroup.POST("/:protocol/token", h.Token)
			authGroup.GET("/:protocol/userinfo", h.UserInfo)
			authGroup.POST("/:protocol/logout", h.Logout)
		}
	}
}
