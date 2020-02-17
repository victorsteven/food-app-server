package interfaces

import (
	"bytes"
	"encoding/json"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	signin func(*entity.User) (map[string]interface{}, map[string]string)
)
type fakeSignin struct {}

func (s *fakeSignin) SignIn(user *entity.User) (map[string]interface{}, map[string]string) {
	return signin(user)
}

func TestSignin_Success(t *testing.T) {
	//Mock all the functions that the function depend on.
	auth.Token = &fakeToken{}
	auth.Auth = &fakeAuth{}
	application.UserApp = &fakeUserApp{}

	getUserEmailPasswordApp = func(*entity.User) (*entity.User, map[string]string) {
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
	createToken  = func(userid uint64) (*auth.TokenDetails, error){
		return &auth.TokenDetails{
			AccessToken:  "this-is-the-access-token",
			RefreshToken: "this-is-the-refresh-token",
			TokenUuid:    "dfsdf-342-34-23-4234-234",
			RefreshUuid:  "sfd-3234-sdfew-34234-df3",
			AtExpires:    12345,
			RtExpires:    1234555,
		}, nil
	}
	createAuth = func(uint64, *auth.TokenDetails) error {
		return nil
	}

	user := &entity.User{
		FirstName: "victor",
		LastName:  "steven",
	}
	details, err := Sign.SignIn(user)
	assert.Nil(t, err)
	assert.EqualValues(t, details["access_token"], "this-is-the-access-token")
	assert.EqualValues(t, details["refresh_token"], "this-is-the-refresh-token")
	assert.EqualValues(t, details["first_name"], "victor")
	assert.EqualValues(t, details["last_name"], "steven")
}


//We dont need to mock the application layer, because we won't get there. So we will use table test to cover all validation errors
func Test_Login_Invalid_Data(t *testing.T) {
	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			//empty email
			inputJSON:  `{"email": "","password": "password"}`,
			statusCode: 422,
		},
		{
			//empty password
			inputJSON:  `{"email": "steven@example.com","password": ""}`,
			statusCode: 422,
		},
		{
			//invalid email
			inputJSON:  `{"email": "stevenexample.com","password": ""}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		r := gin.Default()
		r.POST("/login", Login)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(v.inputJSON))
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
		if validationErr["invalid_email"] != "" {
			assert.Equal(t, validationErr["invalid_email"], "please provide a valid email")
		}
		if validationErr["password_required"] != "" {
			assert.Equal(t, validationErr["password_required"], "password is required")
		}
	}
}


func Test_Login_Success(t *testing.T) {
	//Mock the signin method and return the response:
	Sign = &fakeSignin{} //this where the swap happens
	signin = func(user *entity.User) (map[string]interface{}, map[string]string){
		return map[string]interface{}{
			"access_token": "this-is-the-access-token",
			"refresh_token": "this-is-the-refresh-token",
			"first_name": "victor",
			"last_name": "steven",
			"user_id": 1,
		}, nil
	}
	inputJSON :=  `{"email": "steven@example.com","password": "password"}`
	r := gin.Default()
	r.POST("/login", Login)
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(inputJSON))
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	response := make(map[string]interface{})

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error unmarshalling error %s\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.EqualValues(t, response["access_token"], "this-is-the-access-token")
	assert.EqualValues(t, response["refresh_token"], "this-is-the-refresh-token")
	assert.EqualValues(t, response["first_name"], "victor")
	assert.EqualValues(t, response["last_name"], "steven")
	assert.EqualValues(t, response["user_id"], 1)
}