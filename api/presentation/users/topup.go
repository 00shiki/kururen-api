package users

type TopUpRequest struct {
	Amount float64 `json:"amount" validate:"required"`
}
