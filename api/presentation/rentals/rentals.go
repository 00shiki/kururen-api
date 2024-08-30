package rentals

import "kururen/api/presentation/cars"

type RentalHistoryResponse struct {
	ID             uint               `json:"id"`
	Cars           []cars.CarResponse `json:"cars"`
	PaymentAmount  float64            `json:"payment_amount"`
	PaymentInvoice string             `json:"payment_invoice"`
	StartDate      string             `json:"start_date"`
	EndDate        string             `json:"end_date"`
}
