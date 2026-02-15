package api

import (
	"context"
	"metadata-platform/internal/module/user/model"
	"metadata-platform/internal/module/user/service"
	"metadata-platform/internal/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// SsoUserHandler 用户处理器结构体
type SsoUserHandler struct {
	*utils.BaseHandler
	userService service.SsoUserService
}

// NewSsoUserHandler 创建用户处理器实例
func NewSsoUserHandler(userService service.SsoUserService) *SsoUserHandler {
	return &SsoUserHandler{
		BaseHandler: utils.NewBaseHandler(),
		userService: userService,
	}
}

// SsoCreateUserRequest 创建用户请求结构
type SsoCreateUserRequest struct {
	Account  string `json:"account" form:"account" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Name     string `json:"name" form:"name"`
	Mobile   string `json:"mobile" form:"mobile"`
	Email    string `json:"email" form:"email"`
	Kind     int    `json:"kind" form:"kind"`
	Status   int    `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
}

// SsoUpdateUserRequest 更新用户请求结构
type SsoUpdateUserRequest struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name"`
	Mobile   string `json:"mobile" form:"mobile"`
	Email    string `json:"email" form:"email"`
	Kind     int    `json:"kind" form:"kind"`
	Status   int    `json:"status" form:"status"`
	Remark   string `json:"remark" form:"remark"`
}

// CreateUser 创建用户
func (h *SsoUserHandler) CreateUser(c context.Context, ctx *app.RequestContext) {
	var req SsoCreateUserRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 创建用户模型
	user := &model.SsoUser{
		Account:  req.Account,
		Password: req.Password,
		Name:     req.Name,
		Mobile:   req.Mobile,
		Email:    req.Email,
		Kind:     req.Kind,
		Status:   req.Status,
		Remark:   req.Remark,
	}

	// 调用服务层创建用户
	if err := h.userService.CreateUser(user); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 不返回密码
	user.Password = ""

	ctx.JSON(201, user)
}

// GetUserByID 根据ID获取用户
func (h *SsoUserHandler) GetUserByID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层获取用户
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户不存在"})
		return
	}

	// 不返回密码
	user.Password = ""

	ctx.JSON(200, user)
}

// UpdateUser 更新用户
func (h *SsoUserHandler) UpdateUser(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req SsoUpdateUserRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 调用服务层获取用户
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(404, map[string]string{"error": "用户不存在"})
		return
	}

	// 更新用户字段
	if req.Account != "" {
		user.Account = req.Account
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Mobile != "" {
		user.Mobile = req.Mobile
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Kind != 0 {
		user.Kind = req.Kind
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	if req.Remark != "" {
		user.Remark = req.Remark
	}

	// 调用服务层更新用户
	if err := h.userService.UpdateUser(user); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	// 不返回密码
	user.Password = ""

	ctx.JSON(200, user)
}

// DeleteUser 删除用户
func (h *SsoUserHandler) DeleteUser(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")

	// 调用服务层删除用户
	if err := h.userService.DeleteUser(id); err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(200, map[string]string{"message": "用户删除成功"})
}

// GetAllUsers 获取所有用户
func (h *SsoUserHandler) GetAllUsers(c context.Context, ctx *app.RequestContext) {
	// 调用服务层获取所有用户
	users, err := h.userService.GetAllUsers()
	if err != nil {
		ctx.JSON(500, map[string]string{"error": "获取用户列表失败"})
		return
	}

	// 不返回密码
	for i := range users {
		users[i].Password = ""
	}

	ctx.JSON(200, users)
}
