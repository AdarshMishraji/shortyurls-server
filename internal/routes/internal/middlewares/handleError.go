package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HandleError(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			ctx.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"message": fiber.ErrInternalServerError.Message,
			})
		}
	}()
	return ctx.Next()
}
