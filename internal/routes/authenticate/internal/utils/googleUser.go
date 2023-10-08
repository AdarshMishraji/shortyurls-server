package utils

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

func GetGoogleUser(idToken string, ctx context.Context) (*auth.UserRecord, error) {
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
