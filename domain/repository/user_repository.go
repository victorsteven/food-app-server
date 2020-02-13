package repository

import (
	"food-app/domain/entity"
)

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}


//type DuckRepository interface {
//	Get(ID int64) (*entity.Duck, error)
//	GetAll() ([]*entity.Duck, error)
//	Save(duck *entity.Duck) error
//}

//var userRepository UserRepository
//
//// GetDuckRepository returns the DuckRepository
//func GetUserRepository() UserRepository {
//	return userRepository
//}
//
//// InitDuckRepository injects DuckRepository with its implementation
//func InitUserRepository(r UserRepository) {
//	userRepository = r
//}