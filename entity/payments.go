package entity

type Payment struct {
	ID         uint `gorm:"primary_key"`
	Type       string
	InvoiceURL string
	Amount     float64
}
