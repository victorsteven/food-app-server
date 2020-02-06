package infrastructure

import (
	"errors"
	"fmt"
	"food-app/domain/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)



type repositoryUsersCRUD struct {
	db *gorm.DB
}

func NewRepositoryUsersCRUD(db *gorm.DB) *repositoryUsersCRUD {
	return &repositoryUsersCRUD{db}
}

//var (
//	UserRepo repository.UserRepository = &userRepo{}
//)

//func NewServer(db *gorm.DB) repository.UserRepository {
//	return &Server{DB: db}
//}

func (r *repositoryUsersCRUD) SaveUser(user *entity.User) (*entity.User, error) {
	fmt.Println("WE ENTERED THE INFRASTRUCTURE LAYER")
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repositoryUsersCRUD) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *repositoryUsersCRUD) GetUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

