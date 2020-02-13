package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	userAppSave func(*entity.User) (*entity.User, map[string]string)
)

type fakeApp struct {
}

func (fa fakeApp) GetUsers() ([]entity.User, error) {
	panic("implement me")
}

func (fa fakeApp) GetUser(uint64) (*entity.User, error) {
	panic("implement me")
}

func (fa fakeApp) GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string) {
	panic("implement me")
}

func (fa fakeApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return userAppSave(user)
}

func TestSaveUser_Success(t *testing.T) {
	application.UserApp = &fakeApp{}
	userAppSave = func(*entity.User) (*entity.User, map[string]string) {
		//remember we are running sensitive info such as email and password
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
	r := gin.Default()
	r.POST("/users", SaveUser)
	inputJSON := `{
		"first_name": "victor",
		"last_name": "steven",
		"email": "steven@example.com",
		"password": "password"
	}`
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	user := &entity.User{}

	err = json.Unmarshal(rr.Body.Bytes(), &user)

	fmt.Println("the string: ", string(rr.Body.Bytes()))

	assert.Equal(t, rr.Code, 201)
	assert.EqualValues(t, user.FirstName, "victor")
	assert.EqualValues(t, user.LastName, "steven")
}

//When instead a string an integer is supplied, When attempting to unmarshal input to the user struct, it will fail
func Test_SaveUser_Invalid_Input(t *testing.T) {
	//consider instead of passing a string for  first_name, we passed an integer
	inputJSON :=  `{"first_name": 1234, "last_name": "steven","email": "steven@example.com","password": "password"}`

	r := gin.Default()
	r.POST("/users", SaveUser)
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	invalidInput :=  ""

	err = json.Unmarshal(rr.Body.Bytes(), &invalidInput)
	if err != nil {
		t.Errorf("error unmarshalling error %s\n", err)
	}
	assert.Equal(t, rr.Code, 422)
	assert.Equal(t, invalidInput, "invalid json")
}

//We dont need to mock the application layer, because we won't get there. So we will use table test to cover all validation errors
func Test_SaveUser_Invalidating_Data(t *testing.T) {
	samples := []struct {
		inputJSON  string
		statusCode int
		username   string
		email      string
	}{
		{
			inputJSON:  `{"first_name": "", "last_name": "steven","email": "steven@example.com","password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"first_name": "victor", "last_name": "","email": "steven@example.com","password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"first_name": "victor", "last_name": "steven","email": "","password": "password"}`,
			statusCode: 422,
		},
		{
			inputJSON:  `{"first_name": "victor", "last_name": "steven","email": "steven@example.com","password": ""}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		r := gin.Default()
		r.POST("/users", SaveUser)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		validationErr := make(map[string]string)

		err = json.Unmarshal(rr.Body.Bytes(), &validationErr)
		if err != nil {
			t.Errorf("error unmarshalling error %s\n", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if validationErr["email_required"] != "" {
			assert.Equal(t, validationErr["email_required"], "email is required")
		}
		if validationErr["firstname_required"] != "" {
			assert.Equal(t, validationErr["firstname_required"], "first name is required")
		}
		if validationErr["lastname_required"] != "" {
			assert.Equal(t, validationErr["lastname_required"], "last name is required")
		}
		if validationErr["password_required"] != "" {
			assert.Equal(t, validationErr["password_required"], "password is required")
		}
	}
}

//One of such db error is invalid email, it return that from the application and test.
func TestSaveUser_DB_Error(t *testing.T) {
	application.UserApp = &fakeApp{}
	userAppSave = func(*entity.User) (*entity.User, map[string]string) {
		return nil, map[string]string{
			"email_taken": "email already taken",
		}
	}
	r := gin.Default()
	r.POST("/users", SaveUser)
	inputJSON := `{
		"first_name": "victor",
		"last_name": "steven",
		"email": "steven@example.com",
		"password": "password"
	}`
	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	dbErr := make(map[string]string)
	err = json.Unmarshal(rr.Body.Bytes(), &dbErr)
	if err != nil {
		t.Errorf("cannot unmarshall payload to errMap: %s\n", err)
	}
	assert.Equal(t, rr.Code, 500)
	assert.EqualValues(t, dbErr["email_taken"], "email already taken")
}
