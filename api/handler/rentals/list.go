package rentals

import (
	"github.com/labstack/echo/v4"
	CARS_PRESENTATION "kururen/api/presentation/cars"
	RENTALS_PRESENTATION "kururen/api/presentation/rentals"
	"kururen/pkg/response"
	"net/http"
)

// List godoc
// @Summary List Rental Histories
// @Description Get Rental History List
// @Tags Rentals
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]rentals.RentalHistoryResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security BearerAuth
// @Router /rentals [get]
func (handler *Controller) List(c echo.Context) error {
	userID, ok := Get(c, "user_id")
	if !ok {
		result := response.Response{
			Code:    http.StatusUnauthorized,
			Message: "Invalid token",
		}
		return response.HandleResponse(c, &result)
	}

	rentalHistories, errRentals := handler.rr.GetUserRentalHistories(userID)
	if errRentals != nil {
		result := response.Response{
			Code:    http.StatusInternalServerError,
			Message: errRentals.Error(),
		}
		return response.HandleResponse(c, &result)
	}

	rentalHistoriesResponse := make([]RENTALS_PRESENTATION.RentalHistoryResponse, len(rentalHistories))
	for i, rentalHistory := range rentalHistories {
		cars := make([]CARS_PRESENTATION.CarResponse, len(rentalHistory.Cars))
		for j, car := range rentalHistory.Cars {
			cars[j] = CARS_PRESENTATION.CarResponse{
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
		rentalHistoriesResponse[i] = RENTALS_PRESENTATION.RentalHistoryResponse{
			ID:             rentalHistory.ID,
			Cars:           cars,
			PaymentAmount:  rentalHistory.Payment.Amount,
			PaymentInvoice: rentalHistory.Payment.InvoiceURL,
			StartDate:      rentalHistory.StartDate.Format("2006-01-02"),
			EndDate:        rentalHistory.EndDate.Format("2006-01-02"),
		}
	}
	result := response.Response{
		Code:    http.StatusOK,
		Message: "Success",
	}
	return response.HandleResponse(c, &result)
}
