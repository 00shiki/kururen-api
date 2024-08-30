package entity

import "time"

type User struct {
	ID            uint      `gorm:"primary_key"`
	Name          string    `gorm:"column:name;not null"`
	Email         string    `gorm:"column:email;unique;not null"`
	Username      string    `gorm:"column:username;unique;not null"`
	Password      string    `gorm:"column:password;not null"`
	JwtToken      string    `gorm:"column:jwt_token"`
	DepositAmount float64   `gorm:"column:deposit_amount"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}
