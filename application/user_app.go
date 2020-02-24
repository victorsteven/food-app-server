package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type UserApp struct {
	us repository.UserRepository
}

//var UserApp UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}


//func (u *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
//	return infrastructure.UserRepo.SaveUser(user)
//}
//
//func (u *userApp) GetUser(userId uint64) (*entity.User, error) {
//	return infrastructure.UserRepo.GetUser(userId)
//}
//
//func (u *userApp) GetUsers() ([]entity.User, error) {
//	return infrastructure.UserRepo.GetUsers()
//}
//
//func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
//	return infrastructure.UserRepo.GetUserByEmailAndPassword(user)
//}


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
