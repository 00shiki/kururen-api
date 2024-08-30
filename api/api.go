package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
	"kururen/api/middlewares"
	API_V1 "kururen/api/routes/v1"
	_ "kururen/docs"
	"kururen/pkg/mail"
	PKG_VALIDATOR "kururen/pkg/validator"
	"kururen/pkg/xendit"
	"kururen/repository/cars"
	"kururen/repository/rentals"
	"kururen/repository/users"
)

func Init(e *echo.Echo, db *gorm.DB) {
	log := logrus.New()

	e.Validator = &PKG_VALIDATOR.CustomValidator{Validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middlewares.LogrusMiddleware(log))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	usersRepo := users.NewUsersRepository(db)
	carsRepo := cars.NewCarsRepository(db)
	rentalsRepo := rentals.NewRentalsRepository(db)

	xenditService := xendit.NewXenditService()
	mailService := mail.NewMailService()

	router(
		e,
		usersRepo,
		carsRepo,
		rentalsRepo,
		xenditService,
		mailService,
	)
}

func router(
	e *echo.Echo,
	users users.Repository,
	cars cars.Repository,
	rentals rentals.Repository,
	xeditService xendit.Service,
	mailService mail.Service,
) {
	r := e.Group("/api/v1")
	API_V1.RouteUsers(r, users, mailService)
	API_V1.RouteCars(r, cars)
	API_V1.RouteRentals(r, rentals, users, cars, xeditService, mailService)
}
