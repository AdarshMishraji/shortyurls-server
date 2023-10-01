package url

import (
	"context"
	"errors"
	"shorty-urls-server/database"

	uuid "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func updatePassword(
	urlId string,
	userId string,
	password string,
	ctx context.Context,
) error {
	if password == "" {
		return errors.New("Password cannot be empty")
	}

	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return errors.New("Unable to set ")
	} else {
		if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
			ID:        uuid.MustParse(urlId),
			UserID:    uuid.MustParse(userId),
			IsDeleted: false,
		}).Updates(map[string]interface{}{
			"password": hashedPassword,
		}); tx.Error != nil {
			return errors.New("Unable to update password")
		} else {
			if tx.RowsAffected == 0 {
				return errors.New("Not able to update password")
			}
		}
	}

	return nil
}
