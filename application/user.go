package application

import (
	"food-app/database"
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
)

type UserImpl struct {
	//Repository repository.UserRepository
}
func UserApp() UserAppInterface {
	return &UserImpl{}
}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, error)
	GetUsers() ([]entity.User, error)
}

//func NewUserImpl() userAppInterface {
//	return &UserImpl{}
//}

//type userAppInterface interface {
//	SaveUser(*entity.User) (*entity.User, error)
//	GetUser(id uint64) (*entity.User, error)
//}

//GetUser returns a user
func (u *UserImpl) GetUser(id uint64) (*entity.User, error) {
	//return u.Repository.GetUser(id)
	return nil, nil
}

func (u *UserImpl) SaveUser(user *entity.User) (*entity.User, error) {
	db := database.NewDB()
	conn := infrastructure.NewRepositoryUsersCRUD(db)
	//u, err := entity.User{}
	return conn.SaveUser(user)
}

func (u *UserImpl) GetUsers() ([]entity.User, error) {
	db := database.NewDB()
	conn := infrastructure.NewRepositoryUsersCRUD(db)
	//u, err := entity.User{}
	return conn.GetUsers()
}

