package cars

import "kururen/entity"

func (cr *CarsRepository) UpdateCarStatus(carID uint, newStatus string) error {
	return cr.db.Model(
		&entity.Car{},
	).Where(
		"id = ?",
		carID,
	).Update(
		"availability",
		newStatus,
	).Error
}
