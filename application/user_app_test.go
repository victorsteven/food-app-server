package application

import (
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
	"github.com/stretchr/testify/assert"
	"testing"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS

var (
	saveUserRepo func(*entity.User) (*entity.User, map[string]string)
)

type fakeUserRepo struct {}

func (u *fakeUserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return saveUserRepo(user)
}

func (u *fakeUserRepo) GetUser(uint64) (*entity.User, error) {
	panic("implement me")
}

func (u *fakeUserRepo) GetUsers() ([]entity.User, error) {
	panic("implement me")
}

func (u *fakeUserRepo) GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string) {
	panic("implement me")
}

func TestSaveUser_Success(t *testing.T) {
	infrastructure.UserRepo = &fakeUserRepo{}
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