package handler

import (
	"github.com/gofiber/fiber/v2"
	"go-api/common"
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

// ListRepos godoc
// @Summary List Repository
// @Description List Repository
// @Accept json
// @Produce json
// @Router /api/repositories [Get]
func ListRepos(c *fiber.Ctx) error {
	return common.RespOK(c, "")
}
