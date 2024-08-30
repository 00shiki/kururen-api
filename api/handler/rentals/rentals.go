package rentals

import (
	"kururen/pkg/mail"
	"kururen/pkg/xendit"
	"kururen/repository/cars"
	"kururen/repository/rentals"
	"kururen/repository/users"
	"time"
)

type Controller struct {
	rr rentals.Repository
	ur users.Repository
	cr cars.Repository
	xs xendit.Service
	ms mail.Service
}

func NewHandler(
	rr rentals.Repository,
	ur users.Repository,
	cr cars.Repository,
	xs xendit.Service,
	ms mail.Service,
) *Controller {
	return &Controller{
		rr: rr,
		ur: ur,
		cr: cr,
		xs: xs,
		ms: ms,
	}
}

var (
	Get       = GetID
	Bind      = BindPayload
	Validate  = ValidatePayload
	ParseTime = time.Parse
)
