package v1

import (
	"github.com/labstack/echo/v4"
	"kururen/api/handler/users"
	"kururen/api/middlewares"
	"kururen/pkg/mail"
	USERS_REPO "kururen/repository/users"
)

func RouteUsers(r *echo.Group, ur USERS_REPO.Repository, ms mail.Service) {
	router := r.Group("/users")
	handler := users.NewHandler(ur, ms)
	router.POST("/login", handler.Login)
	router.POST("/register", handler.Register)
	protected := router.Group("")
	protected.Use(middlewares.AuthMiddleware)
	protected.POST("/topup", handler.TopUp)
	protected.GET("/me", handler.Detail)
}
