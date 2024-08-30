package users

import (
	"golang.org/x/crypto/bcrypt"
	"kururen/pkg/mail"
	"kururen/repository/users"
)

type Controller struct {
	ur users.Repository
	ms mail.Service
}

func NewHandler(ur users.Repository, ms mail.Service) *Controller {
	return &Controller{
		ur: ur,
		ms: ms,
	}
}

var (
	Bind             = BindPayload
	Validate         = ValidatePayload
	Get              = GetID
	ComparePassword  = bcrypt.CompareHashAndPassword
	GeneratePassword = bcrypt.GenerateFromPassword
	SignedString     = TokenSignedString
)
