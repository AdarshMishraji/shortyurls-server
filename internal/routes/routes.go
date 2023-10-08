package routes

import (
	"os"
	"shorty-urls-server/internal/routes/authenticate"
	"shorty-urls-server/internal/routes/internal/middlewares"
	"shorty-urls-server/internal/routes/internal/session"
	"shorty-urls-server/internal/routes/redirect"
	"shorty-urls-server/internal/routes/url"
	"shorty-urls-server/internal/routes/urls"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/mustache/v2"
)

func HandleRoutes() {
	port := os.Getenv("PORT")
	engine := mustache.New("./views", ".mustache")

	session.InitSession()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(middlewares.HandleError)

	app.Use(compress.New())

	app.Static("/statics", "./statics", fiber.Static{
		Compress: true,
	})

	app.Use(middlewares.TraceLocation)
	app.Use(middlewares.TraceDevice)

	app.Post("/authenticate", authenticate.Authenticate)

	app.Route("/urls", func(router fiber.Router) {
		router.Use(middlewares.ValidateToken)
		router.Get("/", urls.GetAllUrls)  // get all urls
		router.Get("/meta", urls.GetMeta) // get all urls meta data
	})

	app.Route("/url", func(router fiber.Router) {
		router.Use(middlewares.ValidateToken)

		router.Route("/:urlId", func(router fiber.Router) { // manage url
			router.Get("/", url.GetURLDetails)         // get url data with statistics
			router.Delete("/", url.DeleteShortenedURL) // delete url

			router.Put("/alias", url.UpdateAlias)                    // change alias
			router.Put("/status", url.UpdateStatus)                  // change status
			router.Put("/expiration_time", url.UpdateExpirationTime) // change expiration time

			router.Route("/password", func(router fiber.Router) {
				router.Put("/", url.UpdatePassword)    // update password
				router.Delete("/", url.RemovePassword) // remove password
			})
		})

		router.Post("/", url.GenerateShortenURL) // create shorten url
	})

	app.Post("/password-check", redirect.PasswordCheck)
	app.Get("/:urlAlias", redirect.Redirect)

	app.Listen("127.0.0.1:" + port)
}
