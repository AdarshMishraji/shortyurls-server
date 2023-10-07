package url

import (
	"context"
	"shorty-urls-server/internal/database"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
)

func deleteShortenedURL(urlId string, userId string, ctx context.Context) *fiber.Error {
	userUUID := uuid.MustParse(userId)
	urlUUID := uuid.MustParse(urlId)
	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        &urlUUID,
		UserID:    &userUUID,
		IsActive:  true,
		IsDeleted: false,
	}).Updates(map[string]interface{}{
		"is_deleted": true,
		"is_active":  false,
	}); tx.Error != nil {
		return fiber.ErrInternalServerError
	} else {
		if tx.RowsAffected == 0 {
			return fiber.ErrNotFound
		}
	}

	return nil
}
