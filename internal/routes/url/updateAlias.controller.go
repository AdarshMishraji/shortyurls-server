package url

import (
	"context"
	"shorty-urls-server/internal/database"
	"strings"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
)

func updateAlias(
	urlId string,
	userId string,
	alias string,
	ctx context.Context,
) *fiber.Error {
	if alias == "" {
		return fiber.ErrBadRequest
	}
	userUUID := uuid.MustParse(userId)
	urlUUID := uuid.MustParse(urlId)

	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        &urlUUID,
		UserID:    &userUUID,
		IsDeleted: false,
	}).Updates(map[string]interface{}{
		"alias": alias,
	}); tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "violates unique constraint") {
			return fiber.ErrConflict
		}
		return fiber.ErrInternalServerError
	} else {
		if tx.RowsAffected == 0 {
			return fiber.ErrNotFound
		}
	}

	return nil
}
