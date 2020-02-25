package application

import (
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type UserApp struct {
	us repository.UserRepository
}

var Mine UserAppInterface = &UserApp{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *UserApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	fmt.Println("We Entered the user app")
	return u.us.SaveUser(user)
}

func (u *UserApp) GetUser(userId uint64) (*entity.User, error) {
	fmt.Println("We Entered the user app")
	return u.us.GetUser(userId)
}

func (u *UserApp) GetUsers() ([]entity.User, error) {
	fmt.Println("We Entered the user app")
	return u.us.GetUsers()
}

func (u *UserApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	fmt.Println("We Entered the user app")
	return u.us.GetUserByEmailAndPassword(user)
}
