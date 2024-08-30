package users

import "kururen/entity"

func (ur *UsersRepository) GetUserByUsername(username string) (*entity.User, error) {
	user := new(entity.User)
	if err := ur.db.First(user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UsersRepository) GetUserByID(userID uint) (*entity.User, error) {
	user := new(entity.User)
	if err := ur.db.First(user, userID).Error; err != nil {
		return nil, err
	}
	return user, nil
}
