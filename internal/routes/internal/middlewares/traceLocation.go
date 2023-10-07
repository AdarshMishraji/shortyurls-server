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
	if cachedLocation == nil {
		ip := ctx.IP()
		location := utils.SetLocationInfoToContext(ctx, &ip, nil)
		session.Set("location", location)
		if err := session.Save(); err != nil {
			return err
		}
	} else {
		utils.SetLocationInfoToContext(ctx, nil, cachedLocation.(*utils.LocationInfo))
	}

	return ctx.Next()
}
