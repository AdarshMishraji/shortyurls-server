package url

import (
	"context"
	"encoding/json"
	"os"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func generateShortenURL(url string, userId string, ctx context.Context) (string, *fiber.Error) {
	if url == "" {
		return "", fiber.ErrBadRequest
	}

	alreadyExistsAlias := database.ShortenURL{}

	if err := database.DB.Select("alias").Where(&database.ShortenURL{
		OriginalURL: &url,
	}).First(&alreadyExistsAlias).Error; err == nil {
		if alreadyExistsAlias.Alias != nil && *(alreadyExistsAlias.Alias) != "" {
			shortenedUrl := os.Getenv("SELF_URL") + "/" + *(alreadyExistsAlias.Alias)
			return shortenedUrl, nil
		}

		return "", nil
	} else {
		if err.Error() == "record not found" {
			alias := uuid.New().String()[:8]

			meta, err := utils.GetMetaData(url)
			if err != nil {
				return "", fiber.ErrBadRequest
			}

			metaData, err := json.Marshal(meta)
			if err != nil {
				return "", fiber.ErrInternalServerError
			}
			userUUID := uuid.MustParse(userId)
			shortenURLDataToInsert := database.ShortenURL{
				OriginalURL: &url,
				Alias:       &alias,
				UserID:      &userUUID,
				Meta:        metaData,
			}

			if err := database.DB.WithContext(ctx).Clauses(clause.Returning{
				Columns: []clause.Column{{Name: "alias"}},
			}).Create(&shortenURLDataToInsert).Scan(&shortenURLDataToInsert).Error; err != nil {
				return "", fiber.ErrInternalServerError
			} else {
				shortenedUrl := os.Getenv("SELF_URL") + "/" + *(shortenURLDataToInsert.Alias)
				return shortenedUrl, nil
			}
		} else {
			return "", fiber.ErrInternalServerError
		}
	}
}
