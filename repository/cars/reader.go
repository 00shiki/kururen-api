package cars

import "kururen/entity"

func (cr *CarsRepository) GetCars() ([]entity.Car, error) {
	var cars []entity.Car
	if err := cr.db.Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

func (cr *CarsRepository) GetCarByID(carID uint) (*entity.Car, error) {
	var car *entity.Car
	if err := cr.db.First(&car, carID).Error; err != nil {
		return nil, err
	}
	return car, nil
}
