package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type UserApp struct {
	us repository.UserRepository
}

//UserApp implements the UserAppInterface
var _ UserAppInterface = &UserApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *UserApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *UserApp) GetUser(userId uint64) (*entity.User, error) {
	return u.us.GetUser(userId)
}

func (u *UserApp) GetUsers() ([]entity.User, error) {
	return u.us.GetUsers()
}

func (u *UserApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}
