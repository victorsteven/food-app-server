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

//type repositoryUsersCRUD struct {
//	db *gorm.DB
//}

//func NewUserRepository(db *gorm.DB) repository.UserRepository {
//	return &repositoryUsersCRUD{db}
//}
var UserRepo repository.UserRepository = &Server{}

func (db *Server) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	fmt.Println("WE ENTERED")
	//mine := NewDB()
	dbErr := map[string]string{}
	err := db.DB.Debug().Create(&user).Error
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

func (db *Server) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := db.DB.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (db *Server) GetUsers() ([]entity.User, error) {
	var users []entity.User
	//err := r.db.Debug().Find(&users).Error
	err := db.DB.Debug().Find(&users).Error
	//db.Preload("Orders").Find(&users)
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (db *Server) GetUserByEmailAndPassword(u *entity.User) (*entity.User, map[string]string) {
	var user entity.User
	dbErr := map[string]string{}
	err := db.DB.Debug().Where("email = ?", u.Email).Take(&user).Error
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
