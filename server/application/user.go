package application

import (
	"food-app/server/domain/entity"
	"food-app/server/domain/repository"
)

type UserImpl struct {
	Repository repository.UserRepository
}

func NewUserImpl() userAppInterface {
	return &UserImpl{}
}

type userAppInterface interface {
	SaveUser(*entity.User) (*entity.User, error)
	GetUser(id uint64) (*entity.User, error)
}

func (u *UserImpl) GetUser(id uint64) (*entity.User, error) {
	return u.Repository.GetUser(id)
}

func (u *UserImpl) SaveUser(user *entity.User) (*entity.User, error) {
	return u.Repository.SaveUser(user)
}

