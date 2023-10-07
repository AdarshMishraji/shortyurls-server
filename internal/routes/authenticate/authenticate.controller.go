package authenticate

import (
	"context"
	"shorty-urls-server/internal/database"
	commonUtil "shorty-urls-server/internal/internal/utils"
	authUtils "shorty-urls-server/internal/routes/authenticate/internal/utils"
	"shorty-urls-server/internal/routes/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func authenticate(idToken string, ctx context.Context) (string, *fiber.Error) {
	user, err := authUtils.GetGoogleUser(idToken, ctx)

	if err != nil {
		return "", fiber.ErrUnauthorized
	}

	var tokenString string

	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		provider := "google"
		userToInsert := database.User{
			Email:     &(user.Email),
			Name:      &(user.DisplayName),
			Picture:   &(user.PhotoURL),
			Provider:  &provider,
			IsDeleted: false,
		}

		if err := database.DB.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
			Columns:   []clause.Column{{Name: "email"}},
		}, clause.Returning{
			Columns: []clause.Column{{Name: "id"}},
		}).Create(&userToInsert).Scan(&userToInsert).Error; err != nil {
			return fiber.ErrInternalServerError
		}

		locationInfo := commonUtil.GetLocationInfoFromContext(ctx)
		deviceInfo := commonUtil.GetDeviceInfoFromContext(ctx)

		userLoginHistory := database.UserLoginHistory{
			UserID:   userToInsert.ID,
			Location: datatypes.JSON((&locationInfo).String()),
			Device:   datatypes.JSON((&deviceInfo).String()),
		}

		if err := database.DB.WithContext(ctx).Create(&userLoginHistory).Scan(&userLoginHistory).Error; err != nil {
			return fiber.ErrInternalServerError
		}

		tokenString, err = utils.GenerateJWTToken(
			userToInsert.ID.String(),
			*(userToInsert.Email),
			userLoginHistory.ID.String(),
		)

		if err != nil {
			return fiber.ErrInternalServerError
		}

		return nil
	}); err != nil {
		return "", fiber.ErrInternalServerError
	}

	return tokenString, nil
}
