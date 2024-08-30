package entity

type Car struct {
	ID           uint    `gorm:"primary_key"`
	Model        string  `gorm:"column:model;not null"`
	Brand        string  `gorm:"column:brand;not null"`
	Color        string  `gorm:"column:color"`
	Category     string  `gorm:"column:category"`
	Year         string  `gorm:"column:year"`
	RentalCost   float64 `gorm:"column:rental_cost;not null"`
	Availability string  `gorm:"column:availability"`
}
