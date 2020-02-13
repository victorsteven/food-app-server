package application

import (
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
	"github.com/jinzhu/gorm"
)

type UserImpl struct {
	DB *gorm.DB
}

var UserApp UserAppInterface = &UserImpl{}


type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *UserImpl) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	conn := infrastructure.NewUserRepository(u.DB)
	return conn.SaveUser(user)
	//return infrastructure.UserRepo.SaveUser(user)
}

func (u *UserImpl) GetUser(userId uint64) (*entity.User, error) {
	//return infrastructure.UserRepo.GetUser(userId)
	return nil, nil
}

func (u *UserImpl) GetUsers() ([]entity.User, error) {
	//return infrastructure.UserRepo.GetUsers()
	return nil, nil

}

func (u *UserImpl) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	//return infrastructure.UserRepo.GetUserByEmailAndPassword(user)
	return nil, nil

}
