package utils

import "fmt"

// ErrorCode 错误码类型
type ErrorCode int

// 错误码定义
const (
	// 系统错误
	ErrInternal ErrorCode = 500001 // 内部服务器错误
	ErrDatabase ErrorCode = 500002 // 数据库错误
	ErrCache    ErrorCode = 500003 // 缓存错误
	ErrNetwork  ErrorCode = 500004 // 网络错误

	// 业务错误
	ErrBadRequest    ErrorCode = 400001 // 请求参数错误
	ErrUnauthorized  ErrorCode = 400002 // 未授权
	ErrForbidden     ErrorCode = 400003 // 禁止访问
	ErrNotFound      ErrorCode = 400004 // 资源不存在
	ErrConflict      ErrorCode = 400005 // 资源冲突
	ErrValidation    ErrorCode = 400006 // 数据验证失败
	ErrRateLimit     ErrorCode = 400007 // 超出速率限制

	// 认证错误
	ErrInvalidToken  ErrorCode = 401001 // 无效的令牌
	ErrExpiredToken  ErrorCode = 401002 // 令牌过期
	ErrTokenRequired ErrorCode = 401003 // 需要令牌

	// 第三方服务错误
	ErrThirdPartyAPI ErrorCode = 500101 // 第三方API错误
	ErrThirdPartyAuth ErrorCode = 500102 // 第三方认证错误
)

// AppError 自定义应用错误类型
type AppError struct {
	Code    ErrorCode `json:"code"`    // 错误码
	Message string    `json:"message"` // 错误信息
	Err     error     `json:"-"`       // 原始错误
	Stack   string    `json:"-"`       // 堆栈信息
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 实现errors.Unwrap接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建新的应用错误
func NewAppError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewBadRequestError 创建请求参数错误
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(ErrBadRequest, message, err)
}

// NewUnauthorizedError 创建未授权错误
func NewUnauthorizedError(message string, err error) *AppError {
	return NewAppError(ErrUnauthorized, message, err)
}

// NewForbiddenError 创建禁止访问错误
func NewForbiddenError(message string, err error) *AppError {
	return NewAppError(ErrForbidden, message, err)
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(ErrNotFound, message, err)
}

// NewConflictError 创建资源冲突错误
func NewConflictError(message string, err error) *AppError {
	return NewAppError(ErrConflict, message, err)
}

// NewValidationError 创建数据验证失败错误
func NewValidationError(message string, err error) *AppError {
	return NewAppError(ErrValidation, message, err)
}

// NewInternalError 创建内部服务器错误
func NewInternalError(message string, err error) *AppError {
	return NewAppError(ErrInternal, message, err)
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(message string, err error) *AppError {
	return NewAppError(ErrDatabase, message, err)
}

// NewCacheError 创建缓存错误
func NewCacheError(message string, err error) *AppError {
	return NewAppError(ErrCache, message, err)
}

// NewNetworkError 创建网络错误
func NewNetworkError(message string, err error) *AppError {
	return NewAppError(ErrNetwork, message, err)
}

// NewInvalidTokenError 创建无效的令牌错误
func NewInvalidTokenError(message string, err error) *AppError {
	return NewAppError(ErrInvalidToken, message, err)
}

// NewExpiredTokenError 创建令牌过期错误
func NewExpiredTokenError(message string, err error) *AppError {
	return NewAppError(ErrExpiredToken, message, err)
}

// NewThirdPartyAPIError 创建第三方API错误
func NewThirdPartyAPIError(message string, err error) *AppError {
	return NewAppError(ErrThirdPartyAPI, message, err)
}
