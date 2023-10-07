package utils

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/user_agent"
)

type DeviceInfo struct {
	DeviceType  string `json:"deviceType"`
	OsName      string `json:"osName"`
	BrowserName string `json:"browserName"`
}

func (info *DeviceInfo) String() string {
	return fmt.Sprintf(
		`{"DeviceType":"%s","OsName":"%s","BrowserName":"%s"}`,
		info.DeviceType,
		info.OsName,
		info.BrowserName,
	)
}

func GetDeviceInfo(userAgent string) DeviceInfo {
	ua := user_agent.New(userAgent)
	osName := ua.OS()
	platform := ua.Platform()
	browserName, _ := ua.Browser()

	return DeviceInfo{
		DeviceType:  platform,
		OsName:      osName,
		BrowserName: browserName,
	}
}

func GetDeviceInfoFromContext(ctx context.Context) DeviceInfo {
	return ctx.Value("device").(DeviceInfo)
}

func SetDeviceInfoToContext(ctx *fiber.Ctx, userAgent *string, cachedDeviceInfo *DeviceInfo) *DeviceInfo {
	var deviceInfo DeviceInfo

	if (cachedDeviceInfo != nil) && (cachedDeviceInfo.String() == deviceInfo.String()) {
		deviceInfo = *cachedDeviceInfo
	} else {
		deviceInfo = GetDeviceInfo(*userAgent)
	}

	ctx.SetUserContext(context.WithValue(ctx.UserContext(), "device", deviceInfo))
	return &deviceInfo
}
