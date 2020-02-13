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

type repositoryUser struct {
	db *gorm.DB
}

var UserRepo repository.UserRepository = &repositoryUser{}

//return the interface defined in the repository which the "repositoryUser" satisfies
//func NewRepositoryUser(db *gorm.DB) repository.UserRepository {
//	return &repositoryUser{db}
//}
//func NewRepositoryUser(db *gorm.DB) repository.UserRepository {
//	return &DB{db}
//}

func (r *repositoryUser) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	fmt.Println("INFRASTRUCTURE HERE AGAIN YET")
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		//If the email is already taken
		if strings.Contains(err.Error(), "duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

//func (r *repositoryUser) SaveUser(user *entity.User) (*entity.User, map[string]string) {
//	dbErr := map[string]string{}
//	err := db.Debug().Create(&user).Error
//	if err != nil {
//		//If the email is already taken
//		if strings.Contains(err.Error(), "duplicate") {
//			dbErr["email_taken"] = "email already taken"
//			return nil, dbErr
//		}
//		//any other db error
//		dbErr["db_error"] = "database error"
//		return nil, dbErr
//	}
//	return user, nil
//}

func (r *repositoryUser) GetUser(id uint64) (*entity.User, error) {
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

func (r *repositoryUser) GetUsers() ([]entity.User, error) {
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

func (r *repositoryUser) GetUserByEmailAndPassword(u *entity.User) (*entity.User, map[string]string) {
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
