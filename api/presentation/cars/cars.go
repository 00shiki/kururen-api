package cars

type CarResponse struct {
	ID           uint    `json:"id"`
	Model        string  `json:"model"`
	Brand        string  `json:"brand"`
	Color        string  `json:"color"`
	Category     string  `json:"category"`
	Year         string  `json:"year"`
	RentalCost   float64 `json:"rental_cost"`
	Availability string  `json:"availability"`
}
