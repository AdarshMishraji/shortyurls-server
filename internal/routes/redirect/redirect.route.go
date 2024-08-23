package redirect

import (
	"os"
	"shorty-urls-server/internal/routes/internal/middlewares"
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type PasswordCheckBody struct {
	Password  string `json:"password"`
	Signature string `json:"signature"`
}

func Redirect(ctx *fiber.Ctx) error {
	middlewares.TraceLocation(ctx)
	middlewares.TraceDevice(ctx)
	urlAlias := ctx.Params("urlAlias")
	ip := ctx.Get("Cf-Connecting-Ip")
	if url, err := redirect(urlAlias, ip, string(ctx.Context().UserAgent()), ctx.UserContext(), ctx.Render); err != nil {
		if err == fiber.ErrNotFound {
			return ctx.Render("404", fiber.Map{
				"FrontendURL": os.Getenv("FRONTEND_URL"),
			})
		} else {
			return ctx.Render("internal-error", fiber.Map{
				"FrontendURL": os.Getenv("FRONTEND_URL"),
			})
		}
	} else if url != "" {
		return ctx.Redirect(url, fiber.StatusMovedPermanently)
	} else {
		return nil
	}
}

func PasswordCheck(ctx *fiber.Ctx) error {
	body := new(PasswordCheckBody)

	if err := ctx.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}

	if body.Password == "" || body.Signature == "" {
		return fiber.ErrBadRequest
	}
	ip := ctx.Get("Cf-Connecting-Ip")

	if originalUrl, err := passwordCheck(body.Password, body.Signature, ip, string(ctx.Context().UserAgent()), ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Success",
			Data: fiber.Map{
				"original_url": originalUrl,
			},
		}.SendResponse(ctx)
	}
}
