package persistence

import (
	"food-app/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

//SINCE WE ARE SPINNING UP A DATABASE, THE TESTS HERE ARE INTEGRATION TESTS

//YOU CAN TEST METHOD FAILURES IF YOU HAVE TIME, TO IMPROVE COVERAGE.

func TestSaveUser_Success(t *testing.T) {
	conn, err := DBConn()
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

	conn, err := DBConn()
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

func TestGetUser_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the user
	user, err := seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewUserRepository(conn)
	u, getErr := repo.GetUser(user.ID)

	assert.Nil(t, getErr)
	assert.EqualValues(t, u.Email, "steven@example.com")
	assert.EqualValues(t, u.FirstName, "vic")
	assert.EqualValues(t, u.LastName, "stev")
}

func TestGetUsers_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the users
	_, err = seedUsers(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	repo := NewUserRepository(conn)
	users, getErr := repo.GetUsers()

	assert.Nil(t, getErr)
	assert.EqualValues(t, len(users), 2)
}

func TestGetUserByEmailAndPassword_Success(t *testing.T) {
	conn, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	//seed the user
	u, err := seedUser(conn)
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	var user = &entity.User{
		Email:    "steven@example.com",
		Password: "password",
	}
	repo := NewUserRepository(conn)
	u, getErr := repo.GetUserByEmailAndPassword(user)

	assert.Nil(t, getErr)
	assert.EqualValues(t, u.Email, user.Email)
	//Note, the user password from the database should not be equal to a plane password, because that one is hashed
	assert.NotEqual(t, u.Password, user.Password)
}
