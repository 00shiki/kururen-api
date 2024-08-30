package users

import (
	"gorm.io/gorm"
	"kururen/entity"
)

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	GetUserByUsername(string) (*entity.User, error)
	GetUserByID(uint) (*entity.User, error)
}

type Writer interface {
	CreateUser(*entity.User) error
	UpdateUser(*entity.User) error
}

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}
