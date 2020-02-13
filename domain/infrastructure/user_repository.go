package infrastructure

import (
	"errors"
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"food-app/utils/security"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type repositoryUsersCRUD struct {
	db *gorm.DB
}

var UserRepo repository.UserRepository = &repositoryUsersCRUD{}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &repositoryUsersCRUD{db}
}

//var UserRepo repository.UserRepository = &Server{}

func (r *repositoryUsersCRUD) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	fmt.Println("WE ENTERED")
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		//If the email is already taken
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		fmt.Println("error: ", err)
		dbErr["db_error"] = "database error"
		return nil, dbErr
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
	//err := r.db.Debug().Find(&users).Error
	err := r.db.Debug().Find(&users).Error
	//db.Preload("Orders").Find(&users)
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (r *repositoryUsersCRUD) GetUserByEmailAndPassword(u *entity.User) (*entity.User, map[string]string) {
	var user entity.User
	dbErr := map[string]string{}
	err := r.db.Debug().Where("email = ?", u.Email).Take(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		dbErr["no_user"] = "user not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	//Verify the password
	err = security.VerifyPassword(user.Password, u.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		dbErr["incorrect_password"] = "incorrect password"
		return nil, dbErr
	}
	return &user, nil
}
