package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// log request
			logger.WithFields(logrus.Fields{
				"method": c.Request().Method,
				"path":   c.Request().URL.Path,
				"ip":     c.RealIP(),
			}).Info("Request Received")

			// process next middlewares/handler
			err := next(c)

			// log response status
			logger.WithFields(logrus.Fields{
				"status": c.Response().Status,
			}).Info("Response Sent")

			return err
		}
	}
}
