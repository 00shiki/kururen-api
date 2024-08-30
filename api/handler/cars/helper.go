package cars

import "github.com/labstack/echo/v4"

func BindPayload(c echo.Context, payload interface{}) error {
	return c.Bind(payload)
}

func ValidatePayload(c echo.Context, payload interface{}) error {
	return c.Validate(payload)
}
