package url

import (
	commonUtils "shorty-urls-server/internal/internal/utils"
	utils "shorty-urls-server/internal/routes/internal/utils"

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

func getURLRouteInput[T any](ctx *fiber.Ctx, containsURL bool, containsBody bool) (string, *utils.Claims, *T, *fiber.Error) {
	var urlId string
	if containsURL {
		urlId = ctx.Params("urlId")
	}

	user := utils.GetClaimFromContext(ctx)

	var requestBody *T

	if containsBody {
		requestBody = new(T)
		if err := ctx.BodyParser(requestBody); err != nil {
			return "", nil, nil, fiber.ErrBadRequest
		}

		if containsURL && !commonUtils.IsValidUrlId(urlId) {
			return "", nil, nil, fiber.ErrBadRequest
		}
	}

	return urlId, user, requestBody, nil
}

func GenerateShortenURL(ctx *fiber.Ctx) error {
	_, user, requestBody, err := getURLRouteInput[generateShortenURLRequestBody](ctx, false, true)
	if err != nil {
		return err
	}

	if shortenedUrl, err := generateShortenURL(requestBody.Url, user.UserId, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "URL shortened successfully",
			Data: fiber.Map{
				"shortenedUrl": shortenedUrl,
			},
		}.SendResponse(ctx)
	}
}

func DeleteShortenedURL(ctx *fiber.Ctx) error {
	urlId, user, _, err := getURLRouteInput[interface{}](ctx, true, false)
	if err != nil {
		return err
	}

	if err := deleteShortenedURL(urlId, user.UserId, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "URL deleted successfully",
		}.SendResponse(ctx)
	}
}

func UpdateAlias(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updateAliasRequestBody](ctx, true, true)
	if err != nil {
		return err
	}

	if err := updateAlias(urlId, user.UserId, requestBody.Alias, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Alias updated successfully",
		}.SendResponse(ctx)
	}
}

func UpdateStatus(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updateStatusRequestBody](ctx, true, true)
	if err != nil {
		return err
	}

	if err := updateStatus(urlId, user.UserId, requestBody.IsActive, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Status updated successfully",
		}.SendResponse(ctx)
	}
}

func UpdateExpirationTime(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updateExpirationTimeRequestBody](ctx, true, true)
	if err != nil {
		return err
	}

	if err := updateExpirationTime(urlId, user.UserId, requestBody.ExpirationTime, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Expiration time updated successfully",
		}.SendResponse(ctx)
	}
}

func UpdatePassword(ctx *fiber.Ctx) error {
	urlId, user, requestBody, err := getURLRouteInput[updatePasswordRequestBody](ctx, true, true)
	if err != nil {
		return err
	}

	if err := updatePassword(urlId, user.UserId, requestBody.Password, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Password updated successfully",
		}.SendResponse(ctx)
	}
}

func RemovePassword(ctx *fiber.Ctx) error {
	urlId, user, _, err := getURLRouteInput[interface{}](ctx, true, false)
	if err != nil {
		return err
	}

	if err := removePassword(urlId, user.UserId, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "Password removed successfully",
		}.SendResponse(ctx)
	}
}

func GetURLDetails(ctx *fiber.Ctx) error {
	urlId, user, _, err := getURLRouteInput[interface{}](ctx, true, false)
	if err != nil {
		return err
	}

	if urlDetails, err := getUrlDetails(urlId, user.UserId, ctx.UserContext()); err != nil {
		return err
	} else {
		return utils.Response{
			Message: "URL details fetched successfully",
			Data: fiber.Map{
				"urlDetails": urlDetails,
			},
		}.SendResponse(ctx)
	}
}
