package authenticate

import (
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthenticateRequestBody struct {
	IdToken string `json:"idToken"`
}

func Authenticate(ctx *fiber.Ctx) error {
	requestBody := new(AuthenticateRequestBody)
	if err := ctx.BodyParser(requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if token, err := authenticate(requestBody.IdToken, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Successfully authenticated",
			Data: map[string]interface{}{
				"token": token,
			},
		}.SendResponse(ctx)
	}
}
