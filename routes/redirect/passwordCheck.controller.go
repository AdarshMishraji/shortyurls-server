package redirect

import (
	"shorty-urls-server/database"
	"shorty-urls-server/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func passwordCheck(password string, signature string, ip string, userAgent string, ctx *fiber.Ctx) error {
	if payload, err := validateSignature(signature); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		existingURL := new(database.ShortenURL)

		if err := database.DB.Model(new(database.ShortenURL)).Where(&database.ShortenURL{
			ID: uuid.MustParse(payload.ID),
		}).Select("password", "original_url").Scan(&existingURL).Error; err != nil {
			return internalError(ctx)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(existingURL.Password), []byte(password)); err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid password",
			})
		}

		locationInfo := utils.GetLocationInfo(ip)
		deviceInfo := utils.GetDeviceInfo(userAgent)

		shortenURLRecord := database.ShortenURLVisit{
			ShortenURLID: uuid.MustParse(payload.ID),
			Location:     database.JSON(locationInfo.String()),
			Device:       database.JSON(deviceInfo.String()),
		}

		if err := database.DB.Create(&shortenURLRecord).Error; err != nil {
			return internalError(ctx)
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"url": existingURL.OriginalURL,
		})
	}
}
