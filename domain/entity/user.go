package entity

import (
	"food-app/utils/app_errors"
	"food-app/utils/security"
	"github.com/badoux/checkmail"
	"html"
	"net/http"
	"strings"
	"time"
)

type User struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	FirstName string     `gorm:"size:100;not null;" json:"first_name"`
	LastName  string     `gorm:"size:100;not null;" json:"last_name"`
	Email     string     `gorm:"size:100;not null;unique" json:"email"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

//BeforeSave is a gorm hook
func (u *User) BeforeSave() error {
	hashPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}

func (u *User) Prepare() {
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) *app_errors.UserError {
	var userErr  = app_errors.UserError{}
	var err error

	switch strings.ToLower(action) {
	//case "update":
	//	if u.Email == "" {
	//		err = errors.New("Required Email")
	//		errorMessages["Required_email"] = err.Error()
	//	}
	//	if u.Email != "" {
	//		if err = checkmail.ValidateFormat(u.Email); err != nil {
	//			err = errors.New("Invalid Email")
	//			errorMessages["Invalid_email"] = err.Error()
	//		}
	//	}
	//
	//case "login":
	//	if u.Password == "" {
	//		err = errors.New("Required Password")
	//		errorMessages["Required_password"] = err.Error()
	//	}
	//	if u.Email == "" {
	//		err = errors.New("Required Email")
	//		errorMessages["Required_email"] = err.Error()
	//	}
	//	if u.Email != "" {
	//		if err = checkmail.ValidateFormat(u.Email); err != nil {
	//			err = errors.New("Invalid Email")
	//			errorMessages["Invalid_email"] = err.Error()
	//		}
	//	}
	//case "forgotpassword":
	//	if u.Email == "" {
	//		err = app_errors.New("Required Email")
	//		errorMessages["Required_email"] = err.Error()
	//	}
	//	if u.Email != "" {
	//		if err = checkmail.ValidateFormat(u.Email); err != nil {
	//			err = app_errors.New("Invalid Email")
	//			errorMessages["Invalid_email"] = err.Error()
	//		}
	//	}
	default:
		if u.FirstName == "" {
			userErr.FnErr = "first name is required"
		}
		if u.LastName == "" {
			userErr.LnErr = "last name is required"
		}
		//if u.Password == "" {
		//	err = errors.New("Required Password")
		//	errorMessages["Required_password"] = err.Error()
		//}
		//if u.Password != "" && len(u.Password) < 6 {
		//	err = errors.New("Password should be atleast 6 characters")
		//	errorMessages["Invalid_password"] = err.Error()
		//}
		if u.Email == "" {
			userErr.EmailErr = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				userErr.EmailErr = "invalid email"
			}
		}
	}
	empty := app_errors.UserError{}
	if userErr != empty {
		userErr.StatusErr = http.StatusUnprocessableEntity
		return &userErr
	}
	return nil
}


