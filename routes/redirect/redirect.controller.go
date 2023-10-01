package redirect

import (
	"os"
	"shorty-urls-server/database"
	"shorty-urls-server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func internalError(ctx *fiber.Ctx) error {
	ctx.Render("internal-error", fiber.Map{
		"FrontendURL": os.Getenv("FRONTEND_URL"),
	})
	return nil
}

func redirect(
	urlAlias string,
	ip string,
	userAgent string,
	ctx *fiber.Ctx,
) error {
	if urlAlias == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "URL not found",
		})
	}

	existingURL := database.ShortenURL{}
	if err := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		Alias:     urlAlias,
		IsDeleted: false,
	}).Select("id, original_url, password, expiration_time").First(&existingURL).Error; err != nil {
		if err.Error() == "record not found" {
			return ctx.Render("404", fiber.Map{
				"FrontendURL": os.Getenv("FRONTEND_URL"),
			})
		}
		return internalError(ctx)
	} else {
		if !existingURL.ExpirationTime.IsZero() && existingURL.ExpirationTime.Before(time.Now().UTC()) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "URL expired",
			})
		}

		locationInfo := utils.GetLocationInfo(ip)
		deviceInfo := utils.GetDeviceInfo(userAgent)

		if existingURL.Password != "" {
			if encryptedPayload, err := generateSignature(existingURL.ID.String()); err != nil {
				return internalError(ctx)
			} else {
				return ctx.Render("password-entry", fiber.Map{
					"Signature": encryptedPayload,
					"BaseURL":   os.Getenv("SELF_URL"),
				})
			}
		} else {
			shortenURLRecord := database.ShortenURLVisit{
				ShortenURLID: existingURL.ID,
				Location:     database.JSON(locationInfo.String()),
				Device:       database.JSON(deviceInfo.String()),
			}

			if err := database.DB.Create(&shortenURLRecord).Error; err != nil {
				return internalError(ctx)
			}
			return ctx.Redirect(existingURL.OriginalURL, fiber.StatusMovedPermanently)
		}
	}
}
