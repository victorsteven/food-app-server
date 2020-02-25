package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type UserApp struct {
	us repository.UserRepository
}

var _ UserAppInterface = &UserApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *UserApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.SaveUser(user)
}

func (u *UserApp) GetUser(userId uint64) (*entity.User, error) {
	return u.GetUser(userId)
}

func (u *UserApp) GetUsers() ([]entity.User, error) {
	return u.GetUsers()
}

func (u *UserApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.GetUserByEmailAndPassword(user)
}
