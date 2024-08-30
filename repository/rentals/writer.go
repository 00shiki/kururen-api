package rentals

import "kururen/entity"

func (rr *RentalsRepository) CreateRentalHistory(rentalHistory *entity.RentalHistory) error {
	return rr.db.Create(rentalHistory).Error
}
