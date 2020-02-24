package infrastructure

import (
	"errors"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"food-app/utils/security"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"strings"
)


func NewUserService(db *gorm.DB) repository.UserRepository {
	return &userRepository{db}
}

type userRepository struct {
	db *gorm.DB
}

//The userRepository implements the repository.UserRepository interface
var _ repository.UserRepository = &userRepository{}

//The struct userRepository now implement the UserRepository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	//db := rdbms.NewDB()
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		//If the email is already taken
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return user, nil
}

func (r *userRepository) GetUser(id uint64) (*entity.User, error) {
	//db := rdbms.NewDB()
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

func (r *userRepository) GetUsers() ([]entity.User, error) {
	//db := rdbms.NewDB()
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

func (r *userRepository) GetUserByEmailAndPassword(u *entity.User) (*entity.User, map[string]string) {
	//db := rdbms.NewDB()
	//fmt.Println("WE ARE HERE")
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
