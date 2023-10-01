package url

import (
	"errors"
	"shorty-urls-server/utils"

	"github.com/gofiber/fiber/v2"
)

type generateShortenURLRequestBody struct {
	Url string `json:"url"`
}

type updateAliasRequestBody struct {
	Alias string `json:"alias"`
}

type updateStatusRequestBody struct {
	IsActive bool `json:"is_active"`
}

type updateExpirationTimeRequestBody struct {
	ExpirationTime string `json:"expiration_time"`
}

type updatePasswordRequestBody struct {
	Password string `json:"password"`
}

func getURLRouteInput[T any](ctx *fiber.Ctx, containsURL bool, containsBody bool) (string, *utils.Claims, *T, error) {
	var urlId string
	if containsURL {
		urlId = ctx.Params("urlId")
	}

	user := utils.GetClaimFromContext(ctx)

	var requestBody *T

	if containsBody {
		requestBody = new(T)
		if err := ctx.BodyParser(requestBody); err != nil {
			return "", nil, nil, errors.New("Invalid request")
		}

		if containsURL && !utils.IsValidUrlId(urlId) {
			return "", nil, nil, errors.New("Invalid URL ID")
		}
	}

	return urlId, user, requestBody, nil
}

func GenerateShortenURL(ctx *fiber.Ctx) error {
	_, user, requestBody, err := getURLRouteInput[generateShortenURLRequestBody](ctx, false, true)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if shortenedUrl, err := generateShortenURL(requestBody.Url, user.UserId, ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":      "URL shortened successfully",
			"shortenedUrl": shortenedUrl,
		})
	}
}

func DeleteShortenedURL(ctx *fiber.Ctx) error {
	urlId, user, _, err := getURLRouteInput[interface{}](ctx, true, false)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := deleteShortenedURL(urlId, user.UserId, ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "URL deleted successfully",
		})
	}
}

func UpdateAlias(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updateAliasRequestBody](ctx, true, true)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := updateAlias(urlId, user.UserId, requestBody.Alias, ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Alias updated successfully",
		})
	}
}

func UpdateStatus(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updateStatusRequestBody](ctx, true, true)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := updateStatus(urlId, user.UserId, requestBody.IsActive, ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Status updated successfully",
		})
	}
}

func UpdateExpirationTime(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updateExpirationTimeRequestBody](ctx, true, true)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := updateExpirationTime(urlId, user.UserId, requestBody.ExpirationTime, ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Expiration Time updated successfully",
		})
	}
}

func UpdatePassword(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updatePasswordRequestBody](ctx, true, true)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := updatePassword(urlId, user.UserId, requestBody.Password, ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Password updated successfully",
		})
	}
}

func RemovePassword(ctx *fiber.Ctx) error {
	urlId, user, _, err := getURLRouteInput[interface{}](ctx, true, false)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := removePassword(urlId, user.UserId, ctx.UserContext()); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	} else {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Password removed successfully",
		})
	}
}
