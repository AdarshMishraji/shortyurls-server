package middlewares

import "github.com/gofiber/fiber/v2"

func HandleError(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"message": fiber.ErrInternalServerError.Message,
			})
		}
	}()
	return ctx.Next()
}
