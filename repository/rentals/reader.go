package rentals

import "kururen/entity"

func (rr *RentalsRepository) GetUserRentalHistories(userID uint) ([]entity.RentalHistory, error) {
	var rentalHistories []entity.RentalHistory
	if err := rr.db.Find(&rentalHistories, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return rentalHistories, nil
}

func (rr *RentalsRepository) GetRentalHistoryByID(rentalHistoryID uint) (*entity.RentalHistory, error) {
	rentalHistory := new(entity.RentalHistory)
	if err := rr.db.First(rentalHistory, rentalHistoryID).Error; err != nil {
		return nil, err
	}
	return rentalHistory, nil
}
