package middlewares

import (
	"fmt"
	"shorty-urls-server/internal/internal/utils"
	"shorty-urls-server/internal/routes/internal/session"

	"github.com/gofiber/fiber/v2"
)

func TraceLocation(ctx *fiber.Ctx) error {
	fmt.Println(ctx.Request().Header.String(), ctx.IP(), ctx.Get("X-Forwarded-For"), ctx.IPs())
	session, err := session.SessionStore.Get(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cachedLocation := session.Get("location")
	ip := ctx.IP()
	if cachedLocation == nil {
		location := utils.SetLocationInfoToContext(ctx, &ip, nil)
		session.Set("location", location)
		if err := session.Save(); err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		cachedLocation := cachedLocation.(utils.LocationInfo)
		utils.SetLocationInfoToContext(ctx, &ip, &cachedLocation)
	}

	return ctx.Next()
}
