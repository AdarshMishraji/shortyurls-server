package redirect

import (
	"github.com/gofiber/fiber/v2"
)

type PasswordCheckBody struct {
	Password  string `json:"password"`
	Signature string `json:"signature"`
}

func Redirect(ctx *fiber.Ctx) error {
	urlAlias := ctx.Params("urlAlias")

	return redirect(urlAlias,
		ctx.IP(),
		string(ctx.Context().UserAgent()),
		ctx)
}

func PasswordCheck(ctx *fiber.Ctx) error {
	body := new(PasswordCheckBody)

	if err := ctx.BodyParser(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if body.Password == "" || body.Signature == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	return passwordCheck(body.Password, body.Signature, ctx.IP(), string(ctx.Context().UserAgent()), ctx)
}
