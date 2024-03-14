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
		// helm repo list
		repositories.Get("", handler.ListRepos)
		// helm repo add
		repositories.Post("/:repositories", handler.AddRepo)
		// helm repo remove
		repositories.Delete("/:repositories", handler.RemoveRepo)
		// helm repo update
		repositories.Put("/:repositories", handler.UpdateRepo)
		// helm search chart list
		repositories.Get("/:repositories/charts", handler.ListRepoCharts)

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
