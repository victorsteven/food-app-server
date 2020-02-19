package application

import (
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
)

type userApp struct {
}

var UserApp UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}


func (u *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return infrastructure.UserRepo.SaveUser(user)
}

func (u *userApp) GetUser(userId uint64) (*entity.User, error) {
	return infrastructure.UserRepo.GetUser(userId)
}

func (u *userApp) GetUsers() ([]entity.User, error) {
	return infrastructure.UserRepo.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return infrastructure.UserRepo.GetUserByEmailAndPassword(user)
}

