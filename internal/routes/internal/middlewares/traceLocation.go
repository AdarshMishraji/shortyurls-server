package middlewares

import (
	"shorty-urls-server/internal/internal/utils"
	"shorty-urls-server/internal/routes/internal/session"

	"github.com/gofiber/fiber/v2"
)

func TraceLocation(ctx *fiber.Ctx) error {
	session, err := session.SessionStore.Get(ctx)
	if err != nil {
		return err
	}
	cachedLocation := session.Get("location")
	ip := ctx.IP()
	if cachedLocation == nil {
		location := utils.SetLocationInfoToContext(ctx, &ip, nil)
		session.Set("location", location)
		if err := session.Save(); err != nil {
			return err
		}
	} else {
		cachedLocation := cachedLocation.(utils.LocationInfo)
		utils.SetLocationInfoToContext(ctx, &ip, &cachedLocation)
	}

	return ctx.Next()
}
