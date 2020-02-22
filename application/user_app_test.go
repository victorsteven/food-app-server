package application

import (
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
	"github.com/stretchr/testify/assert"
	"testing"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS FAILURES

var (
	saveUserRepo func(*entity.User) (*entity.User, map[string]string)
	getUserRepo func(userId uint64) (*entity.User, error)
	getUsersRepo func() ([]entity.User, error)
	getUserEmailAndPasswordRepo func(*entity.User) (*entity.User, map[string]string)
)

type fakeUserRepo struct {}

func (u *fakeUserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return saveUserRepo(user)
}
func (u *fakeUserRepo) GetUser(userId uint64) (*entity.User, error) {
	return getUserRepo(userId)
}
func (u *fakeUserRepo) GetUsers() ([]entity.User, error) {
	return getUsersRepo()
}
func (u *fakeUserRepo) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return getUserEmailAndPasswordRepo(user)
}

func TestSaveUser_Success(t *testing.T) {
	infrastructure.UserRepo = &fakeUserRepo{} //swap the real method with the fake
	//Mock the response coming from the infrastructure
	saveUserRepo = func(user *entity.User) (*entity.User,  map[string]string) {
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
			Email:     "steven@example.com",
			Password:  "password",
		}, nil
	}
	user := &entity.User{
		ID:        1,
		FirstName: "victor",
		LastName:  "steven",
		Email:     "steven@example.com",
		Password:  "password",
	}
	u, err := UserApp.SaveUser(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}

func TestGetUser_Success(t *testing.T) {
	infrastructure.UserRepo = &fakeUserRepo{} //swap the real method with the fake
	//Mock the response coming from the infrastructure
	getUserRepo = func(userId uint64) (*entity.User,  error) {
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
			Email:     "steven@example.com",
			Password:  "password",
		}, nil
	}
	userId := uint64(1)
	u, err := UserApp.GetUser(userId)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}

func TestGetUsers_Success(t *testing.T) {
	infrastructure.UserRepo = &fakeUserRepo{} //swap the real method with the fake
	//Mock the response coming from the infrastructure
	getUsersRepo = func() ([]entity.User, error) {
		return []entity.User{
			{
				ID:        1,
				FirstName: "victor",
				LastName:  "steven",
				Email:     "steven@example.com",
				Password:  "password",
			},
			{
				ID:        2,
				FirstName: "kobe",
				LastName:  "bryant",
				Email:     "kobe@example.com",
				Password:  "password",
			},
		}, nil
	}
	users, err := UserApp.GetUsers()
	assert.Nil(t, err)
	assert.EqualValues(t, len(users), 2)
}

func TestGetUserByEmailAndPassword_Success(t *testing.T) {
	infrastructure.UserRepo = &fakeUserRepo{} //swap the real method with the fake
	//Mock the response coming from the infrastructure
	getUserEmailAndPasswordRepo = func(user *entity.User) (*entity.User,  map[string]string) {
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
			Email:     "steven@example.com",
			Password:  "password",
		}, nil
	}
	user := &entity.User{
		ID:        1,
		FirstName: "victor",
		LastName:  "steven",
		Email:     "steven@example.com",
		Password:  "password",
	}
	u, err := UserApp.GetUserByEmailAndPassword(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}