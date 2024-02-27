package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-api/common"
)

// AddRepos
// @Summary Add Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories [Post]
func AddRepos(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"message": "success",
		"data":    "this is test",
	})
}

// ListRepos
// @Summary List Repository
// @Tags Repository
// @Accept json
// @Produce json
// @Router /api/repositories [Get]
func ListRepos(c *fiber.Ctx) error {
	return common.RespOK(c, "")
}
