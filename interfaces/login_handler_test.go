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


func TestSignin_Success(t *testing.T) {
	application.UserApp = &fakeApp{}
	getUserEmailPasswordApp = func(*entity.User) (*entity.User, map[string]string) {
		//remember we are running sensitive info such as email and password
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
	user := &entity.User{
		FirstName: "victor",
		LastName:  "steven",
	}
	details, err := Signin(user)
	fmt.Println(err)
	fmt.Println(details)
}


//We dont need to mock the application layer, because we won't get there. So we will use table test to cover all validation errors
func Test_Login(t *testing.T) {
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