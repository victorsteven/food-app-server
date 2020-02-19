package infrastructure

import (
	"food-app/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSaveUser_Success(t *testing.T) {
	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = entity.User{}
	user.Email = "steven@example.com"
	user.FirstName = "victoria"
	user.LastName = "steven"
	user.Password = "password"

	repo := NewUserRepository(conn)

	u, saveErr := repo.SaveUser(&user)
	assert.Nil(t, saveErr)
	assert.EqualValues(t, u.Email, "steven@example.com")
	assert.EqualValues(t, u.FirstName, "victoria")
	assert.EqualValues(t, u.LastName, "steven")
	//The pasword is supposed to be hashed, so, it should not the same the one we passed:
	assert.NotEqual(t, u.Password, "password")
}

//Failure can be due to duplicate email, etc
//Here, we will attempt saving a user that is already saved
func TestSaveUser_Failure(t *testing.T) {

	conn, err := Database()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the user
	_, err = seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = entity.User{}
	user.Email = "steven@example.com"
	user.FirstName = "Kedu"
	user.LastName = "Nwanne"
	user.Password = "password"

	repo := NewUserRepository(conn)
	u, saveErr := repo.SaveUser(&user)
	dbMsg := map[string]string{
		"email_taken": "email already taken",
	}
	assert.Nil(t, u)
	assert.EqualValues(t, dbMsg, saveErr)
}
