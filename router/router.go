package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go-api/handler"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	api := app.Group("/api")
	// helm repo
	repositories := api.Group("/repositories")
	{
		// helm repo list
		repositories.Post("", handler.AddRepos)

	}

}
