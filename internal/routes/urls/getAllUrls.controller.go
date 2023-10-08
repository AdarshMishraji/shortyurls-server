package urls

import (
	"context"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func getAllUrls(userId string, withMeta bool, ctx context.Context) ([]utils.UrlDetailsResponse, error) {
	selectString := "su.*"
	if withMeta {
		selectString = "JSON_AGG(JSONB_BUILD_OBJECT('device', suv.device, 'location', suv.\"location\", 'visited_at', suv.created_at)) as visits, su.*"
	}

	tx := database.DB.Table("shorten_urls su").
		Select(selectString).
		Where("su.user_id = ?", userId)
	if withMeta {
		tx = tx.Joins("LEFT JOIN shorten_url_visits suv ON suv.shorten_url_id = su.id").
			Group("su.id")
	}

	rows, err := tx.Rows()
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	defer rows.Close()

	urlDetails := []utils.UrlDetailsResponse{}
	for rows.Next() {
		urlDetail := utils.UrlDetails{}
		if err := tx.ScanRows(rows, &urlDetail); err != nil {
			return nil, fiber.ErrInternalServerError
		}

		response, err := utils.GetUrlDetails(urlDetail, withMeta, ctx)
		if err != nil {
			return nil, err
		}

		urlDetails = append(urlDetails, *response)
	}

	return urlDetails, nil
}
