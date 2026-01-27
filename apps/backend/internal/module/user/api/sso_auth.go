package api

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/mssola/user_agent"

	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"
)

// SsoAuthHandler 认证 API 处理器
type SsoAuthHandler struct {
	authService service.SsoAuthService
}

// NewSsoAuthHandler 创建认证 API 处理器实例
func NewSsoAuthHandler(authService service.SsoAuthService) *SsoAuthHandler {
	return &SsoAuthHandler{authService: authService}
}

// SsoLoginRequest 登录请求
type SsoLoginRequest struct {
	Account     string `json:"account" form:"account" binding:"required"`
	Password    string `json:"password" form:"password" binding:"required"`
	TenantID    uint   `json:"tenant_id" form:"tenant_id"`
	CaptchaID   string `json:"captcha_id" form:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" form:"captcha_code" binding:"required"`
}

// SsoRefreshRequest 刷新令牌请求
type SsoRefreshRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

// Login 登录
func (h *SsoAuthHandler) Login(c context.Context, ctx *app.RequestContext) {
	var req SsoLoginRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 校验验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		ctx.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": "验证码不正确或已过期",
		})
		return
	}

	ip := ctx.ClientIP()
	userAgentStr := string(ctx.UserAgent())
	
	// Parse User-Agent
	ua := user_agent.New(userAgentStr)
	browserName, browserVersion := ua.Browser()
	engineName, engineVersion := ua.Engine()
	
	clientInfo := service.ClientInfo{
		IP:             ip,
		UserAgent:      userAgentStr,
		Browser:        browserName,
		BrowserVersion: browserVersion,
		BrowserEngine:  engineName + " " + engineVersion,
		Language:       string(ctx.Request.Header.Get("Accept-Language")),
		OS:             ua.OSInfo().Name,
		OSVersion:      ua.OSInfo().Version,
		OSArch:         "", // difficult to reliably get from standard UA, keep empty or parse specific tokens if needed
		Device:         "",
		Platform:       ua.Platform(),
		Timezone:       "", // Would need client to send it in header/body
		IPLocation:     "", // Would need GeoIP lookup
	}
	
	// Basic device detection
	if ua.Mobile() {
		clientInfo.DeviceType = "Mobile"
		clientInfo.Device = "Mobile"
	} else {
		clientInfo.DeviceType = "Desktop"
		clientInfo.Device = "PC"
	}
	if ua.Bot() {
		clientInfo.DeviceType = "Bot"
		clientInfo.Device = "Bot"
	}

	accessToken, refreshToken, err := h.authService.Login(req.Account, req.Password, req.TenantID, clientInfo)
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]any{
			"code":    401,
			"message": "登录失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]any{
		"code":    200,
		"message": "登录成功",
		"data": map[string]any{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	})
}

// GetCaptcha 获取图形验证码
func (h *SsoAuthHandler) GetCaptcha(c context.Context, ctx *app.RequestContext) {
	id, b64s, err := utils.GenerateCaptcha()
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": "生成验证码失败",
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]any{
		"code":    200,
		"message": "success",
		"data": map[string]any{
			"captcha_id": id,
			"pic_path":   b64s,
		},
	})
}

// Refresh 刷新令牌
func (h *SsoAuthHandler) Refresh(c context.Context, ctx *app.RequestContext) {
	var req SsoRefreshRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	newAccessToken, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]any{
			"code":    401,
			"message": "刷新令牌失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]any{
		"code":    200,
		"message": "刷新成功",
		"data": map[string]any{
			"access_token": newAccessToken,
		},
	})
}

// Logout 退出登录
func (h *SsoAuthHandler) Logout(c context.Context, ctx *app.RequestContext) {
	// Extract info for audit log
	ip := ctx.ClientIP()
	userAgentStr := string(ctx.UserAgent())
	
	// Parse User-Agent
	ua := user_agent.New(userAgentStr)
	browserName, browserVersion := ua.Browser()
	engineName, engineVersion := ua.Engine()

	clientInfo := service.ClientInfo{
		IP:             ip,
		UserAgent:      userAgentStr,
		Browser:        browserName,
		BrowserVersion: browserVersion,
		BrowserEngine:  engineName + " " + engineVersion,
		Language:       string(ctx.Request.Header.Get("Accept-Language")),
		OS:             ua.OSInfo().Name,
		OSVersion:      ua.OSInfo().Version,
		OSArch:         "",
		Device:         "",
		Platform:       ua.Platform(),
	}
	
	if ua.Mobile() {
		clientInfo.DeviceType = "Mobile"
		clientInfo.Device = "Mobile"
	} else {
		clientInfo.DeviceType = "Desktop"
		clientInfo.Device = "PC"
	}

	// UserID should be available if route is authenticated
	userID := ""
	if v, exists := ctx.Get("user_id"); exists {
		userID = v.(string)
	}

	err := h.authService.Logout(c, userID, clientInfo)
	if err != nil {
		ctx.JSON(consts.StatusOK, map[string]any{
			"code":    500,
			"message": "退出失败: " + err.Error(),
		})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]any{
		"code":    200,
		"message": "退出成功",
	})
}

// GetProfile 获取用户信息
func (h *SsoAuthHandler) GetProfile(c context.Context, ctx *app.RequestContext) {
	userID, _ := ctx.Get("user_id")
	if userID == nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]any{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	user, err := h.authService.GetUserInfo(userID.(string))
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]any{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]any{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
	NewPassword string `json:"new_password" form:"new_password" binding:"required"`
}

// ChangePassword 修改密码
func (h *SsoAuthHandler) ChangePassword(c context.Context, ctx *app.RequestContext) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, map[string]any{
			"code":    401,
			"message": "未登录",
		})
		return
	}

	var req ChangePasswordRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	if err := h.authService.ChangePassword(userId.(string), req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]any{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]any{
		"code":    200,
		"message": "密码修改成功",
	})
}
