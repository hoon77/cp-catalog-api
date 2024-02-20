package router

import (
	"github.com/gofiber/fiber/v2"
	"go-api/handler"
)

func APIRoutes(app *fiber.App) {
	api := app.Group("/api")

	// repositories
	repositories := api.Group("/repositories")
	{
		repositories.Post("", handler.AddRepos)
		repositories.Get("", handler.ListRepos)
	}

	// releases
	releases := api.Group("/clusters/:clusterId/namespaces/:namespace/releases")
	{
		// helm list
		releases.Get("", handler.ListReleases)
		// helm get
		releases.Get("/:release", handler.GetReleaseInfo)
		//helm install
		releases.Post("/:release", handler.InstallRelease)
		// helm uninstall
		releases.Delete("/:release", handler.UninstallRelease)

	}

}
