package rentals

type CreateRequest struct {
	Cars        []CarRequest `json:"cars"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	PaymentType string       `json:"payment_type"`
}

type CarRequest struct {
	CarID uint `json:"car_id"`
}

type CreateResponse struct {
	PaymentAmount float64 `json:"payment_amount"`
	InvoiceURL    string  `json:"invoice_url"`
}
