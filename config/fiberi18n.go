package config

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

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
