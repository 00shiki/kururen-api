package cars

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	CARS_PRESENTATION "kururen/api/presentation/cars"
	"kururen/pkg/response"
	"net/http"
)

// List godoc
// @Summary List Car
// @Description Get Car List
// @Tags Cars
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]cars.CarResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /cars [get]
func (handler *Controller) List(c echo.Context) error {
	cars, errList := handler.cr.GetCars()
	if errList != nil {
		if errors.Is(errList, gorm.ErrRecordNotFound) {
			result := response.Response{
				Code:    http.StatusNotFound,
				Message: "cars not found",
			}
			return response.HandleResponse(c, &result)
		}
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errList.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	carsResponse := make([]CARS_PRESENTATION.CarResponse, len(cars))
	for i, car := range cars {
		carsResponse[i] = CARS_PRESENTATION.CarResponse{
			ID:           car.ID,
			Model:        car.Model,
			Brand:        car.Brand,
			Color:        car.Color,
			Category:     car.Category,
			Year:         car.Year,
			RentalCost:   car.RentalCost,
			Availability: car.Availability,
		}
	}
	result := response.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    carsResponse,
	}
	return response.HandleResponse(c, &result)
}
