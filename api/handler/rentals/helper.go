package rentals

import "github.com/labstack/echo/v4"

func GetID(c echo.Context, key string) (uint, bool) {
	value, ok := c.Get(key).(float64)
	return uint(value), ok
}

func BindPayload(c echo.Context, payload interface{}) error {
	return c.Bind(payload)
}

func ValidatePayload(c echo.Context, payload interface{}) error {
	return c.Validate(payload)
}
