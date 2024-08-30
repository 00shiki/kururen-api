package cars

import (
	"gorm.io/gorm"
	"kururen/entity"
)

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	GetCars() ([]entity.Car, error)
	GetCarByID(uint) (*entity.Car, error)
}

type Writer interface {
	UpdateCarStatus(uint, string) error
}

type CarsRepository struct {
	db *gorm.DB
}

func NewCarsRepository(db *gorm.DB) *CarsRepository {
	return &CarsRepository{
		db: db,
	}
}
