package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string                 `json:"message" json:",omitempty"`
	Data    map[string]interface{} `json:"data" json:",omitempty"`
}

func (r Response) SendResponse(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(r)
}
