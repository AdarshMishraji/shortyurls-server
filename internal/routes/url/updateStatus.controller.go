package url

import (
	"context"
	"shorty-urls-server/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func updateStatus(
	urlId string,
	userId string,
	isActive bool,
	ctx context.Context,
) *fiber.Error {
	userUUID := uuid.MustParse(userId)
	urlUUID := uuid.MustParse(urlId)

	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        &urlUUID,
		UserID:    &userUUID,
		IsDeleted: false,
		IsActive:  !isActive,
	}).Updates(map[string]interface{}{
		"is_active": isActive,
	}); tx.Error != nil {
		return fiber.ErrInternalServerError
	} else {
		if tx.RowsAffected == 0 {
			return fiber.ErrNotFound
		}
	}

	return nil
}
