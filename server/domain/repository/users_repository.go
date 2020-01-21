package repository

import "food-app/server/domain/entity"

type userInterface interface {
	SaveUser(*entity.User) (*entity.User, error)
	GetUser(uint64) (*entity.User, error)
}

type user struct {}

func (u *user) SaveUser(*entity.User) (*entity.User, error) {
	panic("implement me")
}

func (u *user) GetUser(uint64) (*entity.User, error) {
	panic("implement me")
}

func NewUser() userInterface {
	return &user{}
}