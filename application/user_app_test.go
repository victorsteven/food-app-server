package application

import (
	"food-app/domain/entity"
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

var userApp UserAppInterface = &fakeUserRepo{} //this is where the real implementation is swap with our fake implementation


func TestSaveUser_Success(t *testing.T) {
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
	u, err := userApp.SaveUser(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}

func TestGetUser_Success(t *testing.T) {
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
	u, err := userApp.GetUser(userId)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}

func TestGetUsers_Success(t *testing.T) {
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
	users, err := userApp.GetUsers()
	assert.Nil(t, err)
	assert.EqualValues(t, len(users), 2)
}

func TestGetUserByEmailAndPassword_Success(t *testing.T) {
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
	u, err := userApp.GetUserByEmailAndPassword(user)
	assert.Nil(t, err)
	assert.EqualValues(t, u.FirstName, "victor")
	assert.EqualValues(t, u.LastName, "steven")
	assert.EqualValues(t, u.Email, "steven@example.com")
}