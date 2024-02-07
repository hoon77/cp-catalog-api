package common

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ResultStatus struct {
	ResultCode     string      `json:"resultCode"`
	ResultMessage  string      `json:"resultMessage"`
	HttpStatusCode int         `json:"httpStatusCode"`
	DetailMessage  string      `json:"detailMessage"`
	Items          interface{} `json:"items"`
}

func RespOK(c *fiber.Ctx, data interface{}) error {
	resultStatus := ResultStatus{
		ResultCode:     RESULT_STATUS_SUCCESS,
		ResultMessage:  localize(c, "OK"),
		HttpStatusCode: 0,
		DetailMessage:  localize(c, "OK"),
		Items:          data,
	}
	return c.Status(200).JSON(resultStatus)

}

func RespErr(c *fiber.Ctx, err error) error {
	log.Errorf("[Error Reason] : %s", err.Error())
	resultStatus := ResultStatus{
		ResultCode:     RESULT_STATUS_FAIL,
		ResultMessage:  err.Error(),
		HttpStatusCode: 400,
		DetailMessage:  err.Error(),
		Items:          make([]string, 0),
	}
	return c.Status(400).JSON(resultStatus)

}

func localize(c *fiber.Ctx, msg string) string {
	localize_msg, _ := fiberi18n.Localize(c, msg)
	return localize_msg
}
