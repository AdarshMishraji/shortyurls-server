package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId  string `json:"userId"`
	Email   string `json:"email"`
	LoginId string `json:"loginId"`
	jwt.RegisteredClaims
}

func (claims *Claims) String() string {
	return fmt.Sprintf(
		`{"UserId":"%s","Email":"%s","LoginId":"%s","RegisteredClaims":{"ExpiresAt":%d,"IssuedAt":%d,"Issuer":"%s"}}`,
		claims.UserId,
		claims.Email,
		claims.LoginId,
		claims.ExpiresAt.Unix(),
		claims.IssuedAt.Unix(),
		claims.Issuer,
	)
}

func GenerateJWTToken(
	userId string,
	email string,
	loginId string,
) (tokenString string, err error) {
	now := time.Now()
	claims := Claims{
		userId,
		email,
		loginId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "shorty-urls-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWTToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	},
		jwt.WithIssuer("shorty-urls-server"),
	)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}

func GetClaimFromContext(ctx *fiber.Ctx) *Claims {
	return ctx.UserContext().Value("user").(*Claims)
}
