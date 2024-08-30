package v1

import (
	"github.com/labstack/echo/v4"
	"kururen/api/handler/cars"
	"kururen/api/middlewares"
	CARS_REPO "kururen/repository/cars"
)

func RouteCars(r *echo.Group, cr CARS_REPO.Repository) {
	router := r.Group("/cars")
	handler := cars.NewHandler(cr)
	router.Use(middlewares.AuthMiddleware)
	router.GET("", handler.List)
	router.GET("/:id", handler.Detail)
	router.PUT("/:id/availability", handler.Availability)
}
