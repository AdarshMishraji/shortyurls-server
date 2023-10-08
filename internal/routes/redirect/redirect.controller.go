package redirect

import (
	"context"
	"os"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/internal/utils"
	redirectUtils "shorty-urls-server/internal/routes/redirect/internal/utils"

	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func redirect(
	urlAlias string,
	ip string,
	userAgent string,
	ctx context.Context,
	render func(string, interface{}, ...string) error,
) (string, *fiber.Error) {
	if urlAlias == "" {
		return "", fiber.ErrNotFound
	}

	existingURL := database.ShortenURL{}
	if err := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		Alias:     &urlAlias,
		IsDeleted: false,
	}).Select("id, original_url, password, expiration_time").First(&existingURL).Error; err != nil {
		if err.Error() == "record not found" {
			return "", fiber.ErrNotFound
		}
		return "", fiber.ErrInternalServerError
	} else {
		expirationTime := existingURL.ExpirationTime
		if expirationTime != nil && !(*expirationTime).IsZero() && (*expirationTime).Before(time.Now().UTC()) {
			return "", fiber.ErrNotFound
		}
		password := existingURL.Password
		if password != nil && *password != "" {
			if encryptedPayload, err := redirectUtils.GenerateSignature(existingURL.ID.String()); err != nil {
				return "", fiber.ErrInternalServerError
			} else {
				err := render("password-entry", fiber.Map{
					"Signature": encryptedPayload,
					"BaseURL":   os.Getenv("SELF_URL"),
				})
				if err != nil {
					return "", fiber.ErrInternalServerError
				}
				return "", nil
			}
		} else {
			locationInfo := utils.GetLocationInfoFromContext(ctx)
			deviceInfo := utils.GetDeviceInfoFromContext(ctx)

			shortenURLRecord := database.ShortenURLVisit{
				ShortenURLID: existingURL.ID,
				Location:     datatypes.JSON((&locationInfo).String()),
				Device:       datatypes.JSON((&deviceInfo).String()),
			}

			if err := database.DB.Create(&shortenURLRecord).Error; err != nil {
				return "", fiber.ErrInternalServerError
			}
			return *(existingURL.OriginalURL), nil
		}
	}
}
