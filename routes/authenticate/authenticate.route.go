package authenticate

import (
	"github.com/gofiber/fiber/v2"
)

type AuthenticateRequestBody struct {
	IdToken string `json:"idToken"`
}

func Authenticate(ctx *fiber.Ctx) error {
	requestBody := new(AuthenticateRequestBody)
	if err := ctx.BodyParser(requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	if token, err := authenticate(requestBody.IdToken, ctx.IP(), string(ctx.Context().UserAgent()), ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	}
}
