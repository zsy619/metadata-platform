package utils

import (
	"github.com/mojocn/base64Captcha"
)

// DefaultCaptchaStore 默认验证码存储
var DefaultCaptchaStore = base64Captcha.DefaultMemStore

// DefaultCaptchaConfig 默认验证码配置
var DefaultCaptchaConfig = base64Captcha.DriverDigit{
	Height:   80,
	Width:    240,
	Length:   6,
	MaxSkew:  0.7,
	DotCount: 90,
}

// GenerateCaptcha 生成图形验证码
func GenerateCaptcha() (id string, b64s string, err error) {
	driver := &DefaultCaptchaConfig
	c := base64Captcha.NewCaptcha(driver, DefaultCaptchaStore)
	id, b64s, _, err = c.Generate()
	return
}

// VerifyCaptcha 校验验证码
func VerifyCaptcha(id string, answer string) bool {
	return DefaultCaptchaStore.Verify(id, answer, true)
}
