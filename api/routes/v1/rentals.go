package v1

import (
	"github.com/labstack/echo/v4"
	"kururen/api/handler/rentals"
	"kururen/api/middlewares"
	"kururen/pkg/mail"
	"kururen/pkg/xendit"
	CARS_REPO "kururen/repository/cars"
	RENTALS_REPO "kururen/repository/rentals"
	USERS_REPO "kururen/repository/users"
)

func RouteRentals(
	r *echo.Group,
	rr RENTALS_REPO.Repository,
	ur USERS_REPO.Repository,
	cr CARS_REPO.Repository,
	xs xendit.Service,
	ms mail.Service,
) {
	router := r.Group("/rentals")
	handler := rentals.NewHandler(rr, ur, cr, xs, ms)
	router.Use(middlewares.AuthMiddleware)
	router.GET("", handler.List)
	router.POST("", handler.Create)
}
