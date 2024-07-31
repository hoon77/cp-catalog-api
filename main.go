package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-api/config"
	_ "go-api/docs"
	"go-api/handler"
	"go-api/middleware"
	"go-api/router"
	"helm.sh/helm/v3/pkg/repo"
	"os"
)

func init() {
	config.InitEnvConfigs()
}

// @title Container Platform Helm Rest API
// @version 1.0
// @description K-PaaS Container Platform Helm Rest API
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	log.Info("Hello, Helm Rest API!")
	makeRepoConfig()
	app := fiber.New()
	middleware.SetupLocalize(app)
	middleware.FiberMiddleware(app)
	router.SwaggerRoute(app)
	router.APIRoutes(app)
	handler.Settings()
	err := app.Listen(config.Env.ServerPort)
	if err != nil {
		log.Fatal("Server is not running! Reason: %v", err)
	}
}

func makeRepoConfig() {
	// Check repositories.yaml exists
	if _, err := os.Stat(config.Env.HelmRepoConfig); os.IsNotExist(err) {
		log.Info("repositories.yaml does not exist...")
		repositories := repo.NewFile()
		log.Infof("Create repositories.yaml...(path : %s)", config.Env.HelmRepoConfig)
		if err = repositories.WriteFile(config.Env.HelmRepoConfig, 0600); err != nil {
			log.Infof("Failed to create repositories.yaml ...%s", err)
		}
	}

	// Check repository cache path exists
	if err := os.MkdirAll(config.Env.HelmRepoCache, os.ModePerm); err != nil {
		log.Infof("Failed to create cache directory(path : %s)...%s", config.Env.HelmRepoCache, err)
	}

	// Check repository ca file path exists
	if err := os.MkdirAll(config.Env.HelmRepoCA, os.ModePerm); err != nil {
		log.Infof("Failed to create ca directory(path : %s)...%s", config.Env.HelmRepoCA, err)
	}
}
