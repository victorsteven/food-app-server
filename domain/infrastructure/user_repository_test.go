package infrastructure

import (
	"food-app/domain/entity"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

//This test involve database interaction.
//So, we will test for success and failure due to duplicate email
//Take note that is not the layer we validate the input, what we are strictly interested here is hitting the database
func TestUserRepo_SaveUser(t *testing.T) {
	//refresh the database:
	err := refreshUserTable()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	sampleData := []struct{
		user *entity.User
		ErrMsg map[string]string
		StatusCode int
	}{
		{
			//Everything goes well
			user: &entity.User{
				FirstName: "Kobe",
				LastName: "Byrant",
				Email: 	"kobe@example.com",
				Password: "password",
			},
			ErrMsg: nil,
			StatusCode: http.StatusCreated,
		},
		{
			//An attempt to register with the same email
			user: &entity.User{
				FirstName: "Michael",
				LastName: "Jordan",
				Email: 	"kobe@example.com",
				Password: "password",
			},
			ErrMsg: map[string]string{
				"email_taken": "email already taken",
			},
			StatusCode: http.StatusInternalServerError,
		},
	}
	for _, v := range sampleData {
		repo := NewRepositoryUser(server.db)
		u, saveErr := repo.SaveUser(v.user)

		assert.EqualValues(t, saveErr, v.ErrMsg)

		if v.StatusCode == http.StatusCreated {
			assert.EqualValues(t, u, v.user)
		} else { //status of 500, etc
			dupEmail := map[string]string{"email_taken": "email already taken"}
			assert.EqualValues(t,  dupEmail, v.ErrMsg)
		}
	}
}
