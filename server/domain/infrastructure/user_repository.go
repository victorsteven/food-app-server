package infrastructure

import (
	"errors"
	"food-app/server/domain/entity"
	"food-app/server/domain/repository"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"github.com/jinzhu/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepository {
	return &userRepo{db: db}
}

func (u *userRepo) SaveUser(user *entity.User) (*entity.User, error) {
	err := u.db.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := u.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

