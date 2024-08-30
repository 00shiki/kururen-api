package users

type DetailResponse struct {
	Username      string  `json:"username"`
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	DepositAmount float64 `json:"deposit_amount"`
}
