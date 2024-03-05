package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go-api/config"
	"go-api/handler"
	"go-api/middleware"
	"go-api/router"
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
	app := fiber.New()
	middleware.FiberMiddleware(app)
	middleware.SetupLocalize(app)
	router.SwaggerRoute(app)
	router.APIRoutes(app)
	handler.Settings()
	err := app.Listen(config.Env.ServerPort)
	if err != nil {
		log.Fatal("Server is not running! Reason: %v", err)
	}
}
