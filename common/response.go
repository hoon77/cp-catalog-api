package common

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
)

type ResultStatus struct {
	ResultCode     string `json:"resultCode"`
	ResultMessage  string `json:"resultMessage"`
	HttpStatusCode int    `json:"httpStatusCode"`
	DetailMessage  string `json:"detailMessage"`
}

// func respOK(c *fiber.Ctx, data interface{})
func RespOK(c *fiber.Ctx) error {
	resultStatus := ResultStatus{
		ResultCode:     localize(c, "OK"),
		ResultMessage:  "message",
		HttpStatusCode: 0,
		DetailMessage:  "detail Message",
	}
	return c.Status(200).JSON(resultStatus)

}

func localize(c *fiber.Ctx, msg string) string {
	localize_msg, _ := fiberi18n.Localize(c, msg)
	return localize_msg
}
