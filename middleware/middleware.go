package middleware

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go-api/config"
	"golang.org/x/text/language"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(
		// Add CORS to each route.
		cors.New(),
		// Add simple logger.
		logger.New(),
		// Add basic auth
		basicauth.New(basicauth.Config{
			Users: map[string]string{
				config.Env.AuthUserName: config.Env.AuthPassword,
			},
		}),
	)
}

func SetupLocalize(app *fiber.App) {
	app.Use(
		fiberi18n.New(&fiberi18n.Config{
			RootPath:         "./localize",
			AcceptLanguages:  []language.Tag{language.English, language.Korean},
			DefaultLanguage:  language.English,
			FormatBundleFile: "json",
		}),
	)
}
