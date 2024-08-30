package cars

type AvailabilityRequest struct {
	Status string `json:"status" validate:"required"`
}
