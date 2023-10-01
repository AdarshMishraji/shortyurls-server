package url

import (
	"context"
	"errors"
	"os"
	"shorty-urls-server/database"

	uuid "github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func generateShortenURL(url string, userId string, ctx context.Context) (string, error) {
	if url == "" {
		return "", errors.New("Invalid request")
	}

	alreadyExistsAlias := database.ShortenURL{}

	if err := database.DB.Select("alias").Where(&database.ShortenURL{
		OriginalURL: url,
	}).First(&alreadyExistsAlias).Error; err == nil {
		if alreadyExistsAlias.Alias != "" {
			shortenedUrl := os.Getenv("SELF_URL") + "/" + alreadyExistsAlias.Alias
			return shortenedUrl, nil
		}

		return "", nil
	} else {
		if err.Error() == "record not found" {
			alias := uuid.New().String()[:8]

			shortenURLDataToInsert := database.ShortenURL{
				OriginalURL: url,
				Alias:       alias,
				UserID:      uuid.MustParse(userId),
			}

			if err := database.DB.WithContext(ctx).Clauses(clause.Returning{
				Columns: []clause.Column{{Name: "alias"}},
			}).Create(&shortenURLDataToInsert).Scan(&shortenURLDataToInsert).Error; err != nil {
				return "", errors.New("Unable to create shortened URL")
			} else {
				shortenedUrl := os.Getenv("SELF_URL") + "/" + shortenURLDataToInsert.Alias
				return shortenedUrl, nil
			}
		} else {
			return "", errors.New("Unable to create shortened URL")
		}
	}
}
