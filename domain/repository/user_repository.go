package repository

import (
	"food-app/domain/entity"
)

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
}



