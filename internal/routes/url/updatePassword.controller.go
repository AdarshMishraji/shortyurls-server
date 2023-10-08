package url

import (
	"context"
	"shorty-urls-server/internal/database"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func updatePassword(
	urlId string,
	userId string,
	password string,
	ctx context.Context,
) *fiber.Error {
	if password == "" {
		return fiber.ErrBadRequest
	}

	if len(password) < 8 {
		return fiber.ErrBadRequest
	}

	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return fiber.ErrInternalServerError
	} else {
		userUUID := uuid.MustParse(userId)
		urlUUID := uuid.MustParse(urlId)

		if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
			ID:        &urlUUID,
			UserID:    &userUUID,
			IsDeleted: false,
		}).Updates(map[string]interface{}{
			"password": hashedPassword,
		}); tx.Error != nil {
			return fiber.ErrInternalServerError
		} else {
			if tx.RowsAffected == 0 {
				return fiber.ErrNotFound
			}
		}
	}

	return nil
}
