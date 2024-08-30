package response

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type Response struct {
	Code    int         `json:"-"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HandleResponse(c echo.Context, res *Response) error {
	if res.Code == http.StatusInternalServerError {
		log.Print(res.Message)
		res.Message = http.StatusText(res.Code)
	}
	return c.JSON(res.Code, &res)
}
