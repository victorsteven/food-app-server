package application

import (
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
)

type UserImpl struct {
	//Repository repository.UserRepository
}

var UserApp UserAppInterface = &UserImpl{}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

//GetUser returns a user
//func (u *UserImpl) GetUser(id uint64) (*entity.User, error) {
//	//return u.Repository.GetUser(id)
//	return nil, nil
//}

func (u *UserImpl) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	db := rdbms.NewDB()
	conn := infrastructure.NewUserRepository(db)
	return conn.SaveUser(user)
}

func (u *UserImpl) GetUser(userId uint64) (*entity.User, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewUserRepository(db)
	return conn.GetUser(userId)
}

func (u *UserImpl) GetUsers() ([]entity.User, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewUserRepository(db)
	//u, err := entity.User{}
	return conn.GetUsers()
}

func (u *UserImpl) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	db := rdbms.NewDB()
	conn := infrastructure.NewUserRepository(db)
	foundUser, err := conn.GetUserByEmailAndPassword(user)
	if err != nil {
		return nil, err
	}
	return foundUser, nil
}