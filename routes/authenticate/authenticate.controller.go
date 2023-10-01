package authenticate

import (
	"context"
	"errors"
	"shorty-urls-server/database"
	"shorty-urls-server/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

func authenticate(idToken string, ip string, userAgent string, ctx context.Context) (string, error) {
	user, err := _getGoogleUser(idToken, ctx)

	if err != nil {
		return "", errors.New("Invalid id token")
	}

	var tokenString string

	if err = database.DB.Transaction(func(tx *gorm.DB) error {
		userToInsert := database.User{
			Email:     user.Email,
			Name:      user.DisplayName,
			Picture:   user.PhotoURL,
			Provider:  "google",
			IsDeleted: false,
		}

		if err := database.DB.WithContext(ctx).Clauses(clause.OnConflict{
			DoNothing: true,
			Columns:   []clause.Column{{Name: "email"}},
		}, clause.Returning{
			Columns: []clause.Column{{Name: "id"}},
		}).Create(&userToInsert).Scan(&userToInsert).Error; err != nil {
			return errors.New("Failed to create user")
		}

		locationInfo := utils.GetLocationInfo(ip)
		deviceInfo := utils.GetDeviceInfo(userAgent)

		userLoginHistory := database.UserLoginHistory{
			UserID:   userToInsert.ID,
			Location: database.JSON(locationInfo.String()),
			Device:   database.JSON(deviceInfo.String()),
		}

		if err := database.DB.WithContext(ctx).Create(&userLoginHistory).Scan(&userLoginHistory).Error; err != nil {

			return errors.New("Failed to create user login history")
		}

		tokenString, err = utils.GenerateJWTToken(
			userToInsert.ID.String(),
			userToInsert.Email,
			userLoginHistory.ID.String(),
		)

		if err != nil {
			return errors.New("Failed to generate token")
		}

		return nil
	}); err != nil {
		return "", err
	}

	return tokenString, nil
}

func _getGoogleUser(idToken string, ctx context.Context) (*auth.UserRecord, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	userRecord, err := client.GetUser(ctx, token.UID)
	if err != nil {
		return nil, err
	}

	return userRecord, nil
}
