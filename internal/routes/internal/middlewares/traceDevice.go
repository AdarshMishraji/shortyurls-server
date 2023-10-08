package middlewares

import (
	"shorty-urls-server/internal/internal/utils"
	"shorty-urls-server/internal/routes/internal/session"

	"github.com/gofiber/fiber/v2"
)

func TraceDevice(ctx *fiber.Ctx) error {
	session, err := session.SessionStore.Get(ctx)
	if err != nil {
		return err
	}
	cachedDevice := session.Get("device")
	userAgent := string(ctx.Context().UserAgent())
	if cachedDevice == nil {
		device := utils.SetDeviceInfoToContext(ctx, &userAgent, nil)
		session.Set("device", device)
		if err := session.Save(); err != nil {
			return err
		}
	} else {
		cachedDevice := cachedDevice.(utils.DeviceInfo)
		utils.SetDeviceInfoToContext(ctx, &userAgent, &cachedDevice)
	}

	return ctx.Next()
}
