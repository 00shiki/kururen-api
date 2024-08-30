package cars

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	CARS_PRESENTATION "kururen/api/presentation/cars"
	"kururen/pkg/response"
	"net/http"
)

// Availability godoc
// @Summary Update Car Availability
// @Description Change Car Availability
// @Tags Cars
// @Accept json
// @Produce json
// @Param RequestBody body cars.AvailabilityRequest true "Car Availability Request"
// @Param id path int true "Car ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /cars/{id}/availability [put]
func (handler *Controller) Availability(c echo.Context) error {
	payload := new(CARS_PRESENTATION.AvailabilityRequest)
	if errBind := Bind(c, payload); errBind != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: errBind.Error(),
		}
		return response.HandleResponse(c, &result)
	}
	if errValidate := Validate(c, payload); errValidate != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: errValidate.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	idStr := c.Param("id")
	id, errConv := Atoi(idStr)
	if errConv != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid car ID",
		}
		return response.HandleResponse(c, &result)
	}

	errUpdate := handler.cr.UpdateCarStatus(uint(id), payload.Status)
	if errUpdate != nil {
		if errors.Is(errUpdate, gorm.ErrRecordNotFound) {
			result := response.Response{
				Code:    http.StatusNotFound,
				Message: "Car not found",
			}
			return response.HandleResponse(c, &result)
		}
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errUpdate.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	result := response.Response{
		Code:    http.StatusOK,
		Message: "Success Update",
	}
	return response.HandleResponse(c, &result)
}
