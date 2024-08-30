package cars

import (
	"kururen/repository/cars"
	"strconv"
)

type Controller struct {
	cr cars.Repository
}

func NewHandler(cr cars.Repository) *Controller {
	return &Controller{
		cr: cr,
	}
}

var (
	Atoi     = strconv.Atoi
	Bind     = BindPayload
	Validate = ValidatePayload
)
