package url

import (
	"context"
	"errors"
	"shorty-urls-server/database"
	"strings"

	uuid "github.com/google/uuid"
)

func updateAlias(
	urlId string,
	userId string,
	alias string,
	ctx context.Context,
) error {
	if alias == "" {
		return errors.New("Alias cannot be empty")
	}

	if tx := database.DB.Model(&database.ShortenURL{}).Where(&database.ShortenURL{
		ID:        uuid.MustParse(urlId),
		UserID:    uuid.MustParse(userId),
		IsDeleted: false,
	}).Updates(map[string]interface{}{
		"alias": alias,
	}); tx.Error != nil {
		if strings.Contains(tx.Error.Error(), "violates unique constraint") {
			return errors.New("Alias already exists")
		}
		return errors.New("Unable to update alias")
	} else {
		if tx.RowsAffected == 0 {
			return errors.New("Not able to update alias")
		}
	}

	return nil
}
