package application

import (
	"food-app/database/rdbms"
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
	//db := rdbms.NewDB()
	//conn := infrastructure.NewUserRepository(db)
	//return conn.SaveUser(user)
	return infrastructure.UserRepo.SaveUser(user)
}

func (u *userApp) GetUser(userId uint64) (*entity.User, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewUserRepository(db)
	return conn.GetUser(userId)
}

func (u *userApp) GetUsers() ([]entity.User, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewUserRepository(db)
	return conn.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	db := rdbms.NewDB()
	conn := infrastructure.NewUserRepository(db)
	foundUser, err := conn.GetUserByEmailAndPassword(user)
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}

