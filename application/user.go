package application

import (
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
	"food-app/utils/app_errors"
)

type UserImpl struct {
	//Repository repository.UserRepository
}

func UserApp() UserAppInterface {
	return &UserImpl{}
}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, *app_errors.UserError)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(string, string) (*entity.User, error)
}

//GetUser returns a user
//func (u *UserImpl) GetUser(id uint64) (*entity.User, error) {
//	//return u.Repository.GetUser(id)
//	return nil, nil
//}

func (u *UserImpl) SaveUser(user *entity.User) (*entity.User, *app_errors.UserError) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryUsersCRUD(db)
	return conn.SaveUser(user)
}

func (u *UserImpl) GetUser(userId uint64) (*entity.User, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryUsersCRUD(db)
	return conn.GetUser(userId)
}

func (u *UserImpl) GetUsers() ([]entity.User, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryUsersCRUD(db)
	//u, err := entity.User{}
	return conn.GetUsers()
}

func (u *UserImpl) GetUserByEmailAndPassword(email string, password string) (*entity.User, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryUsersCRUD(db)
	user, err := conn.GetUserByEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
