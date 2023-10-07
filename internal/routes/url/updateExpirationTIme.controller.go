package url

import (
	"context"
	"shorty-urls-server/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
)

func updateExpirationTime(
	urlId string,
	userId string,
	expirationTime string,
	ctx context.Context,
) *fiber.Error {
	if expirationTime == "" {
		return fiber.ErrBadRequest
	}

	date, err := time.Parse(time.RFC3339, expirationTime)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if date.Before(time.Now()) {
		return fiber.ErrBadRequest
	}
	userUUID := uuid.MustParse(userId)
	urlUUID := uuid.MustParse(urlId)

	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        &urlUUID,
		UserID:    &userUUID,
		IsDeleted: false,
	}).Updates(map[string]interface{}{
		"expiration_time": date,
	}); tx.Error != nil {
		return fiber.ErrInternalServerError
	} else {
		if tx.RowsAffected == 0 {
			return fiber.ErrNotFound
		}
	}

	return nil
}
