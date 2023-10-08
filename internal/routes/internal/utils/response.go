package utils

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (r Response) SendResponse(ctx *fiber.Ctx) error {
	if r.Data != nil {
		data, err := json.Marshal(r.Data)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &r.Data)
		r.Data = removeEmptyFields(r.Data)
	}
	return ctx.Status(fiber.StatusOK).JSON(r)
}

// recursively remove empty fields from map or slice
func removeEmptyFields(input interface{}) interface{} {
	switch input.(type) {
	case map[string]interface{}:
		output := make(map[string]interface{})
		for k, v := range input.(map[string]interface{}) {
			if k == "CreatedAt" || k == "UpdatedAt" || k == "DeletedAt" {
				continue
			}
			if v != nil {
				output[k] = removeEmptyFields(v)
			}
		}
		return output
	case []interface{}:
		output := make([]interface{}, 0)
		for _, v := range input.([]interface{}) {
			if v != nil {
				output = append(output, removeEmptyFields(v))
			}
		}
		return output
	default:
		return input
	}
}
