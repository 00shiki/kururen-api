package users

import "kururen/entity"

func (ur *UsersRepository) CreateUser(user *entity.User) error {
	return ur.db.Create(user).Error
}

func (ur *UsersRepository) UpdateUser(user *entity.User) error {
	return ur.db.Save(user).Error
}
