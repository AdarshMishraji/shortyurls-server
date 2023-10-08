package middlewares

import (
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func ValidateToken(ctx *fiber.Ctx) error {
	headers := ctx.GetReqHeaders()

	token := headers["Authorization"]

	if token == "" {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	decodedToken, err := utils.ValidateJWTToken(token)

	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	utils.SetClaimToContext(ctx, decodedToken)

	return ctx.Next()
}
