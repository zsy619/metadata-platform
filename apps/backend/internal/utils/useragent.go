package utils

import (
	"strings"

	"github.com/mssola/user_agent"
)

// ClientInfo 客户端信息
type ClientInfo struct {
	IP               string
	UserAgent        string
	Browser          string
	BrowserVersion   string
	BrowserEngine    string
	Language         string
	OS               string
	OSVersion        string
	OSArch           string
	DeviceType       string
	DeviceModel      string
	Device           string // Keep for backward compatibility or general device string
	Platform         string
	ScreenResolution string
	Timezone         string
	IPLocation       string
}

// ParseUserAgent 解析 User-Agent 字符串,返回客户端信息
// userAgentStr: User-Agent 原始字符串
// clientIP: 客户端 IP 地址
// acceptLanguage: Accept-Language 请求头
func ParseUserAgent(userAgentStr string, clientIP string, acceptLanguage string) ClientInfo {
	ua := user_agent.New(userAgentStr)
	
	browserName, browserVersion := ua.Browser()
	engineName, engineVersion := ua.Engine()
	
	clientInfo := ClientInfo{
		IP:             clientIP,
		UserAgent:      userAgentStr,
		Browser:        browserName,
		BrowserVersion: browserVersion,
		BrowserEngine:  engineName + " " + engineVersion,
		Language:       acceptLanguage,
		OS:             ua.OSInfo().Name,
		OSVersion:      ua.OSInfo().Version,
		OSArch:         parseOSArch(userAgentStr),
		Platform:       ua.Platform(),
		DeviceModel:    parseDeviceModel(userAgentStr, ua.OSInfo().Name),
		ScreenResolution: "",
		Timezone:       "", // Would need client to send it in header/body
		IPLocation:     "", // Would need GeoIP lookup
	}
	
	// Basic device detection
	if ua.Mobile() {
		clientInfo.DeviceType = "Mobile"
		if clientInfo.DeviceModel != "" {
			clientInfo.Device = clientInfo.DeviceModel
		} else {
			clientInfo.Device = "Mobile"
		}
	} else {
		clientInfo.DeviceType = "Desktop"
		clientInfo.Device = "PC"
	}
	
	if ua.Bot() {
		clientInfo.DeviceType = "Bot"
		clientInfo.Device = "Bot"
	}
	
	return clientInfo
}

// 辅助解析操作系统架构
func parseOSArch(ua string) string {
	if strings.Contains(ua, "x86_64") || strings.Contains(ua, "Win64") || strings.Contains(ua, "x64") {
		return "x64"
	}
	if strings.Contains(ua, "arm64") || strings.Contains(ua, "aarch64") {
		return "arm64"
	}
	if strings.Contains(ua, "i686") || strings.Contains(ua, "i386") {
		return "x86"
	}
	if strings.Contains(ua, "arm") {
		return "arm"
	}
	return ""
}

// 辅助解析设备型号
func parseDeviceModel(ua string, os string) string {
	if strings.Contains(os, "Android") {
		// Android UAs usually look like: ...; Android 10; SM-G960F Build/...
		parts := strings.Split(ua, ";")
		for _, part := range parts {
			if strings.Contains(part, "Build/") {
				subParts := strings.Split(strings.TrimSpace(part), " ")
				if len(subParts) > 0 {
					return subParts[0]
				}
			}
		}
	}
	if strings.Contains(ua, "iPhone") {
		return "iPhone"
	}
	if strings.Contains(ua, "iPad") {
		return "iPad"
	}
	return ""
}
