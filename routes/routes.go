package routes

import (
	"os"
	"shorty-urls-server/middlewares"
	"shorty-urls-server/routes/authenticate"
	"shorty-urls-server/routes/redirect"
	"shorty-urls-server/routes/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/template/mustache/v2"
)

func HandleRoutes() {
	port := os.Getenv("PORT")
	engine := mustache.New("./views", ".mustache")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Post("/authenticate", authenticate.Authenticate)
	app.Use(compress.New())

	app.Static("/public", "./public", fiber.Static{
		Compress: true,
	})

	// app.Route("/urls", func(router fiber.Router) {
	// 	router.Use(middlewares.ValidateToken)
	// 	router.Get("/")     // get all urls
	// 	router.Get("/meta") // get all urls meta data
	// })

	app.Route("/url", func(router fiber.Router) {
		router.Use(middlewares.ValidateToken)

		router.Route("/:urlId", func(router fiber.Router) { // manage url
			// 	router.Get("/")    // get url data with statistics
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
