package router

import (
	"github.com/gofiber/fiber/v2"
	"go-api/handler"
)

func APIRoutes(app *fiber.App) {
	api := app.Group("/api")

	repositories := api.Group("/repositories")
	{
		repositories.Post("", handler.AddRepos)
		repositories.Get("", handler.ListRepos)
	}

}
