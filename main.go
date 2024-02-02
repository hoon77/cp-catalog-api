package main

import (
	"github.com/gofiber/fiber/v2"
	_ "go-api/docs"
	"go-api/router"
)

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
	app := fiber.New()

	router.SetupRoutes(app)

	err := app.Listen(":8080")
	if err != nil {
		return
	}
}
