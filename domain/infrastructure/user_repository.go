package infrastructure

import (
	"errors"
	"food-app/domain/entity"
	"food-app/utils/app_errors"
	"food-app/utils/security"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type repositoryUsersCRUD struct {
	db *gorm.DB
}

func NewRepositoryUsersCRUD(db *gorm.DB) *repositoryUsersCRUD {
	return &repositoryUsersCRUD{db}
}

func (r *repositoryUsersCRUD) SaveUser(user *entity.User) (*entity.User, *app_errors.UserError) {
	dbError := &app_errors.UserError{}
	validateErr := user.Validate("")
	if validateErr != nil {
		return nil, validateErr
	}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			dbError.EmailErr = "email already taken"
			dbError.StatusErr = http.StatusInternalServerError
		}
		return nil, dbError
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

func (r *repositoryUsersCRUD) GetUserByEmailAndPassword(email, password string) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	err = security.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, err
	}
	return &user, nil
}
