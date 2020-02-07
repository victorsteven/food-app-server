package repository

import (
	"food-app/domain/entity"
	"food-app/utils/app_errors"
)

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, *app_errors.UserError)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
}
