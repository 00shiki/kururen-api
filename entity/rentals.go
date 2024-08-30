package entity

import "time"

type RentalHistory struct {
	ID         uint         `gorm:"primary_key"`
	StartDate  time.Time    `gorm:"not null"`
	EndDate    time.Time    `gorm:"not null"`
	UserID     uint         `gorm:"not null"`
	PaymentID  uint         `gorm:"not null"`
	User       User         `gorm:"foreignKey:user_id"`
	Payment    Payment      `gorm:"foreignKey:payment_id;references:id"`
	Cars       []*Car       `gorm:"many2many:car_rentals"`
	CarRentals []CarRentals `gorm:"foreignKey:rental_history_id"`
	Duration   float64      `gorm:"-"`
}

func (RentalHistory) TableName() string {
	return "rental_histories"
}

type CarRentals struct {
	ID              uint `gorm:"primary_key"`
	CarID           uint
	RentalHistoryID uint
}
