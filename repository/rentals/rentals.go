package rentals

import (
	"gorm.io/gorm"
	"kururen/entity"
)

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	GetUserRentalHistories(uint) ([]entity.RentalHistory, error)
	GetRentalHistoryByID(uint) (*entity.RentalHistory, error)
}

type Writer interface {
	CreateRentalHistory(*entity.RentalHistory) error
}

type RentalsRepository struct {
	db *gorm.DB
}

func NewRentalsRepository(db *gorm.DB) *RentalsRepository {
	return &RentalsRepository{
		db: db,
	}
}
