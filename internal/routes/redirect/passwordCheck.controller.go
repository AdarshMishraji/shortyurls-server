package redirect

import (
	"context"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/internal/utils"
	redirectUtils "shorty-urls-server/internal/routes/redirect/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

func passwordCheck(password string, signature string, ip string, userAgent string, ctx context.Context) (string, *fiber.Error) {
	if payload, err := redirectUtils.ValidateSignature(signature); err != nil {
		return "", fiber.ErrBadRequest
	} else {
		existingURL := new(database.ShortenURL)
		urlId := uuid.MustParse(payload.ID)
		if err := database.DB.Model(new(database.ShortenURL)).Where(&database.ShortenURL{
			ID: &urlId,
		}).Select("password", "original_url").Scan(&existingURL).Error; err != nil {
			return "", fiber.ErrInternalServerError
		}

		if err := bcrypt.CompareHashAndPassword([]byte(*(existingURL.Password)), []byte(password)); err != nil {
			return "", fiber.ErrUnauthorized
		}
		locationInfo := utils.GetLocationInfoFromContext(ctx)
		deviceInfo := utils.GetLocationInfoFromContext(ctx)
		shortedURLID := uuid.MustParse(payload.ID)
		shortenURLRecord := database.ShortenURLVisit{
			ShortenURLID: &shortedURLID,
			Location:     datatypes.JSON((&locationInfo).String()),
			Device:       datatypes.JSON((&deviceInfo).String()),
		}

		if err := database.DB.Create(&shortenURLRecord).Error; err != nil {
			return "", fiber.ErrInternalServerError
		}
		return *(existingURL.OriginalURL), nil
	}
}
