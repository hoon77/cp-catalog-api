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
		repositories.Post("", handler.AddRepo)
		repositories.Delete("/:repositories", handler.RemoveRepo)
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
		//helm upgrade
		releases.Put("/:release", handler.UpgradeRelease)
		// helm rollback
		releases.Put("/:release/versions/:reversion", handler.RollbackRelease)
		// helm uninstall
		releases.Delete("/:release", handler.UninstallRelease)
		// helm release history
		releases.Get("/:release/histories", handler.GetReleaseHistories)
		// helm release resources status
		releases.Get("/:release/resources", handler.GetReleaseResources)
	}

	charts := api.Group("/repositories/:repositories/charts/:charts")
	{
		charts.Get("/versions", handler.GetChartVersions)
		charts.Get("/info", handler.GetChartInfo)
	}

}
