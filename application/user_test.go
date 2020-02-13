package application

import (
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
	"testing"
	"time"
)

//we will mock the SaveUser method from the repository, so we can achieve unit test on the SaveUser method of the application:

var (
	saveUserRepo func(*entity.User) (*entity.User, map[string]string)
)

type fakeRepo struct {}

func (fr *fakeRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return saveUserRepo(user)
}

func (fr *fakeRepo) GetUser(uint64) (*entity.User, error) {
	panic("implement me")
}

func (fr *fakeRepo) GetUsers() ([]entity.User, error) {
	panic("implement me")
}

func (fr *fakeRepo) GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string) {
	panic("implement me")
}
//
//func (fr *fakeRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
//	return saveUserRepo(user)
//}

//app := UserImpl{DB: server.DB}
//var server = interfaces.Server{}

func TestUserImpl_SaveUser(t *testing.T) {
	infrastructure.UserRepo = &fakeRepo{}
	saveUserRepo = func(user *entity.User) (*entity.User,  map[string]string) {
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
			Email:     "steven@example.com",
			Password:  "password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		}, nil
	}

	user := &entity.User{
		ID:        1,
		FirstName: "victor",
		LastName:  "steven",
		Email:     "steven@example.com",
		Password:  "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
	u, err := UserApp().SaveUser(user)
	fmt.Println(err)
	fmt.Println(u)
	//db := rdbms.NewDB()
	//db.
}

