package url

import (
	"context"
	"shorty-urls-server/internal/database"
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func getUrlDetails(urlId string, userId string, ctx context.Context) (*utils.UrlDetailsResponse, *fiber.Error) {
	urlDetails := utils.UrlDetails{}

	err := database.DB.Table("shorten_urls su").
		Select("JSON_AGG(JSONB_BUILD_OBJECT('device', suv.device, 'location', suv.\"location\", 'visited_at', suv.created_at)) as visits, su.*").
		Joins("LEFT JOIN shorten_url_visits suv ON suv.shorten_url_id = su.id").
		Where("su.id = ? AND su.user_id = ?", urlId, userId).
		Group("su.id").
		Find(&urlDetails).Error

	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	response, err := utils.GetUrlDetails(urlDetails, true, ctx)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	return response, nil
}
