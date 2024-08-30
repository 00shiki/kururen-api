package cars

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	CARS_PRESENTATION "kururen/api/presentation/cars"
	"kururen/pkg/response"
	"net/http"
)

// Detail godoc
// @Summary Detail Car
// @Description Get Car Detail
// @Tags Cars
// @Accept json
// @Produce json
// @Param id path int true "Car ID"
// @Success 200 {object} response.Response{data=cars.CarResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /cars/{id} [get]
func (handler *Controller) Detail(c echo.Context) error {
	idStr := c.Param("id")
	id, errConv := Atoi(idStr)
	if errConv != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid car ID",
		}
		return response.HandleResponse(c, &result)
	}

	car, errCar := handler.cr.GetCarByID(uint(id))
	if errCar != nil {
		if errors.Is(errCar, gorm.ErrRecordNotFound) {
			result := response.Response{
				Code:    http.StatusNotFound,
				Message: "car not found",
			}
			return response.HandleResponse(c, &result)
		}
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errCar.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	result := response.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data: CARS_PRESENTATION.CarResponse{
			ID:           car.ID,
			Model:        car.Model,
			Brand:        car.Brand,
			Color:        car.Color,
			Category:     car.Category,
			Year:         car.Year,
			RentalCost:   car.RentalCost,
			Availability: car.Availability,
		},
	}
	return response.HandleResponse(c, &result)
}
