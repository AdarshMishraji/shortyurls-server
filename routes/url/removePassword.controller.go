package url

import (
	"context"
	"errors"
	"shorty-urls-server/database"

	uuid "github.com/google/uuid"
)

func removePassword(urlId string, userId string, ctx context.Context) error {
	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        uuid.MustParse(urlId),
		UserID:    uuid.MustParse(userId),
		IsDeleted: false,
	}).Updates(map[string]interface{}{
		"password": nil,
	}); tx.Error != nil {
		return errors.New("Unable to remove password")
	} else {
		if tx.RowsAffected == 0 {
			return errors.New("Not able to remove password")
		}
	}

	return nil
}
