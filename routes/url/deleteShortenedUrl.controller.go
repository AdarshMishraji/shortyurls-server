package url

import (
	"context"
	"errors"
	"shorty-urls-server/database"

	uuid "github.com/google/uuid"
)

func deleteShortenedURL(urlId string, userId string, ctx context.Context) error {
	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        uuid.MustParse(urlId),
		UserID:    uuid.MustParse(userId),
		IsActive:  true,
		IsDeleted: false,
	}).Updates(map[string]interface{}{
		"is_deleted": true,
		"is_active":  false,
	}); tx.Error != nil {
		return errors.New("Unable to delete URL")
	} else {
		if tx.RowsAffected == 0 {
			return errors.New("Not able to delete URL")
		}
	}

	return nil
}
