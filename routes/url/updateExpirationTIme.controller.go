package url

import (
	"context"
	"errors"
	"shorty-urls-server/database"
	"time"

	uuid "github.com/google/uuid"
)

func updateExpirationTime(
	urlId string,
	userId string,
	expirationTime string,
	ctx context.Context,
) error {
	if expirationTime == "" {
		return errors.New("Expiration time cannot be empty")
	}

	date, err := time.Parse(time.RFC3339, expirationTime)
	if err != nil {
		return errors.New("Invalid expiration time")
	}

	if date.Before(time.Now()) {
		return errors.New("Expiration time cannot be in the past")
	}

	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        uuid.MustParse(urlId),
		UserID:    uuid.MustParse(userId),
		IsDeleted: false,
	}).Updates(map[string]interface{}{
		"expiration_time": date,
	}); tx.Error != nil {
		return errors.New("Unable to update expiration time")
	} else {
		if tx.RowsAffected == 0 {
			return errors.New("Not able to update expiration time")
		}
	}

	return nil
}
