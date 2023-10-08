package urls

import (
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAllUrls(ctx *fiber.Ctx) error {
	user := utils.GetClaimFromContext(ctx)
	withMeta := ctx.Query("with_meta") == "true"
	println(withMeta)

	if urls, err := getAllUrls(user.UserId, withMeta, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Success",
			Data: fiber.Map{
				"urls": urls,
			},
		}.SendResponse(ctx)
	}
}

func GetMeta(ctx *fiber.Ctx) error {
	user := utils.GetClaimFromContext(ctx)

	if meta, err := getMeta(user.UserId, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Success",
			Data: fiber.Map{
				"meta": meta,
			},
		}.SendResponse(ctx)
	}
}
