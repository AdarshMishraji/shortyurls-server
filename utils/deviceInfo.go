package utils

import (
	"fmt"

	"github.com/mssola/user_agent"
)

type DeviceInfo struct {
	DeviceName  string `json:"deviceName"`
	DeviceType  string `json:"deviceType"`
	OsName      string `json:"osName"`
	BrowserName string `json:"browserName"`
}

func (info *DeviceInfo) String() string {
	return fmt.Sprintf(
		`{"DeviceName":"%s","DeviceType":"%s","OsName":"%s","BrowserName":"%s"}`,
		info.DeviceName,
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
