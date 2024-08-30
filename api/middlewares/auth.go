package middlewares

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"kururen/pkg/response"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" {
			result := response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Missing Authorization header",
			}
			return response.HandleResponse(c, &result)
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			result := response.Response{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}
			return response.HandleResponse(c, &result)
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			result := response.Response{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Token",
			}
			return response.HandleResponse(c, &result)
		}

		c.Set("user_id", claims["user_id"])

		return next(c)
	}
}
