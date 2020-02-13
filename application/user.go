package application

import (
	"fmt"
	//"food-app/database/rdbms"
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
	"github.com/jinzhu/gorm"
)

type UserImpl struct {
	DB *gorm.DB
	//Repository repository.UserRepository
}

func (u *UserImpl) GetUsers() ([]entity.User, error) {
	panic("implement me")
}

func (u *UserImpl) GetUser(uint64) (*entity.User, error) {
	panic("implement me")
}

func (u *UserImpl) GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string) {
	panic("implement me")
}

func UserApp() UserAppInterface {
	return &UserImpl{}
}

//func DB(db *gorm.DB)  *gorm.DB {
//	return UserImpl{db}
//}

type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *UserImpl) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	fmt.Println("THE APPLICATION")
	//db := rdbms.NewDB()
	//conn := infrastructure.NewRepositoryUser(u.DB)
	//return repository.GetUserRepository().SaveUser(user)
	return infrastructure.UserRepo.SaveUser(user)

}

//func (u *UserImpl) GetUser(userId uint64) (*entity.User, error) {
//	db := rdbms.NewDB()
//	conn := infrastructure.NewRepositoryUser(db)
//	return conn.GetUser(userId)
//}
//
//func (u *UserImpl) GetUsers() ([]entity.User, error) {
//	db := rdbms.NewDB()
//	conn := infrastructure.NewRepositoryUser(db)
//	//u, err := entity.User{}
//	return conn.GetUsers()
//}
//
//func (u *UserImpl) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
//	db := rdbms.NewDB()
//	conn := infrastructure.NewRepositoryUser(db)
//	foundUser, err := conn.GetUserByEmailAndPassword(user)
//	if err != nil {
//		return nil, err
//	}
//	return foundUser, nil
//}
