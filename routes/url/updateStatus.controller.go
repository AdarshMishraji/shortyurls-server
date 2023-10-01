package url

import (
	"context"
	"errors"
	"shorty-urls-server/database"

	"github.com/google/uuid"
)

func updateStatus(
	urlId string,
	userId string,
	isActive bool,
	ctx context.Context,
) error {
	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        uuid.MustParse(urlId),
		UserID:    uuid.MustParse(userId),
		IsDeleted: false,
		IsActive:  !isActive,
	}).Updates(map[string]interface{}{
		"is_active": isActive,
	}); tx.Error != nil {
		return errors.New("Unable to update status")
	} else {
		if tx.RowsAffected == 0 {
			return errors.New("Not able to update status")
		}
	}

	return nil
}
