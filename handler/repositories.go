package handler

import (
	"github.com/gofiber/fiber/v2"
)

// AddRepos godoc
// @Summary Add Repository
// @Description Add Repository
// @Accept json
// @Produce json
// @Router /api/repositories [Post]
func AddRepos(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    "this is test",
	})
}
