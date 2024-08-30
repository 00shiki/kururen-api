package rentals

import (
	"fmt"
	"github.com/labstack/echo/v4"
	RENTALS_PRESENTATION "kururen/api/presentation/rentals"
	"kururen/entity"
	"kururen/pkg/response"
	"net/http"
)

// Create godoc
// @Summary Create Rentals
// @Description Create new rental histories
// @Tags Rentals
// @Accept json
// @Produce json
// @Param RequestBody body rentals.CreateRequest true "Rentals Request"
// @Success 201 {object} response.Response{data=rentals.CreateResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /rentals [post]
func (handler *Controller) Create(c echo.Context) error {
	userID, ok := Get(c, "user_id")
	if !ok {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid token",
		}
		return response.HandleResponse(c, &result)
	}

	user, errUser := handler.ur.GetUserByID(userID)
	if errUser != nil {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "user not found",
		}
		return response.HandleResponse(c, &result)
	}

	payload := new(RENTALS_PRESENTATION.CreateRequest)
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

	startDate, errParseStart := ParseTime("2006-01-02", payload.StartDate)
	if errParseStart != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: errParseStart.Error(),
		}
		return response.HandleResponse(c, &result)
	}
	endDate, errParseEnd := ParseTime("2006-01-02", payload.EndDate)
	if errParseEnd != nil {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: errParseEnd.Error(),
		}
		return response.HandleResponse(c, &result)
	}
	if startDate.After(endDate) {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: "start date cannot be after end date",
		}
		return response.HandleResponse(c, &result)
	}

	duration := endDate.Sub(startDate).Hours() / 24

	var totalCost float64
	cars := make([]*entity.Car, len(payload.Cars))
	for i, carID := range payload.Cars {
		car, errCar := handler.cr.GetCarByID(carID.CarID)
		if errCar != nil {
			result := response.Response{
				Code:    http.StatusBadRequest,
				Message: "car not found",
			}
			return response.HandleResponse(c, &result)
		}
		if car.Availability != "available" {
			result := response.Response{
				Code:    http.StatusBadRequest,
				Message: "car not available",
			}
			return response.HandleResponse(c, &result)
		}
		totalCost += car.RentalCost
		cars[i] = car
	}
	if duration*totalCost > user.DepositAmount {
		result := response.Response{
			Code:    http.StatusBadRequest,
			Message: "Insufficient deposit amount",
		}
		return response.HandleResponse(c, &result)
	}

	rentalHistory := entity.RentalHistory{
		User:      *user,
		Cars:      cars,
		StartDate: startDate,
		EndDate:   endDate,
		Duration:  duration,
		Payment: entity.Payment{
			Type:   payload.PaymentType,
			Amount: duration * totalCost,
		},
	}
	errInvoice := handler.xs.CreateInvoice(&rentalHistory)
	if errInvoice != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errInvoice.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	user.DepositAmount -= duration * totalCost
	errUpdateUser := handler.ur.UpdateUser(user)
	if errUpdateUser != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errUpdateUser.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	errSend := handler.ms.SendMail(
		user.Email,
		user.Name,
		"Payment Invoice",
		fmt.Sprintf("<a href='%s'>click this link to finish payment</a>", rentalHistory.Payment.InvoiceURL),
	)
	if errSend != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errSend.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	errCreate := handler.rr.CreateRentalHistory(&rentalHistory)
	if errCreate != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errCreate.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	for _, car := range cars {
		errUpdateCar := handler.cr.UpdateCarStatus(car.ID, "booked")
		if errUpdateCar != nil {
			result := response.Response{
				Code:    http.StatusInternalServerError,
				Message: errUpdateCar.Error(),
			}
			return response.HandleResponse(c, &result)
		}
	}

	result := response.Response{
		Code:    http.StatusCreated,
		Message: "successfully created",
		Data: RENTALS_PRESENTATION.CreateResponse{
			PaymentAmount: rentalHistory.Payment.Amount,
			InvoiceURL:    rentalHistory.Payment.InvoiceURL,
		},
	}
	return response.HandleResponse(c, &result)
}
