package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"food-app/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)


const (
	Dbdriver = "postgres"
	DbHost     = "127.0.0.1"
	DbPort     = "5432"
	DbName   = "food-app"
	DbPassword = "password"
	DbUser     = "steven"

	redis_host = "127.0.0.1"
	redis_port = "6379"
	redis_password = ""
)

var users = &Users{}


func TestMain(m *testing.M) {

	//if err := godotenv.Load(os.ExpandEnv("./../.env")); err != nil {
	//	log.Println("no env gotten")
	//}
	//
	//services, err := infrastructure.NewServices(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName)
	//if err != nil {
	//	panic(err)
	//}
	//defer services.Close()
	//services.Automigrate()


	//redisService, err := auth.NewRedisDB(redis_host, redis_port, redis_password)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tk  := auth.NewToken()
	//
	//users = NewUsers( tk)
	//os.Exit(m.Run())
}

//type fakeUsers struct {
//	us application.UserAppInterface
//	rd auth.AuthInterface
//	tk auth.TokenInterface
//}
//var faker = fakeUsers{}
//var users = NewUsers(faker.us, faker.rd, faker.tk)

//var userHandler UserHandlerInterface = &fakeUserApp{} //this is where the real implementation is swap with our fake implementation



//var users UserHandlerInterface = &fakeUserApp{}

//var foodApp FoodAppInterface = &fakeFoodRepo{} //this is where the real implementation is swap with our fake implementation

//IF YOU HAVE TIME, YOU CAN TEST ALL FAILURE CASES TO IMPROVE COVERAGE

func TestSaveUser_Success(t *testing.T) {
	saveUserApp = func(*entity.User) (*entity.User, map[string]string) {
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
	r := gin.Default()
	r.POST("/users", users.SaveUser)
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

	assert.Equal(t, rr.Code, 201)
	assert.EqualValues(t, user.FirstName, "victor")
	assert.EqualValues(t, user.LastName, "steven")
}

//We dont need to mock the application layer, because we won't get there. So we will use table test to cover all validation errors
func Test_SaveUser_Invalidating_Data(t *testing.T) {
	samples := []struct {
		inputJSON  string
		statusCode int
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
		{
			//invalid email
			inputJSON:  `{"email": "stevenexample.com","password": ""}`,
			statusCode: 422,
		},
		{
			//When instead a string an integer is supplied, When attempting to unmarshal input to the user struct, it will fail
			inputJSON: `{"first_name": 1234, "last_name": "steven","email": "steven@example.com","password": "password"}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		r := gin.Default()
		r.POST("/users", users.SaveUser)
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
		fmt.Println("validator error: ", validationErr)
		assert.Equal(t, rr.Code, v.statusCode)

		if validationErr["email_required"] != "" {
			assert.Equal(t, validationErr["email_required"], "email is required")
		}
		if validationErr["invalid_email"] != "" {
			assert.Equal(t, validationErr["invalid_email"], "please provide a valid email")
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
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}

//One of such db error is invalid email, it return that from the application and test.
func TestSaveUser_DB_Error(t *testing.T) {
	//application.UserApp = &fakeUserApp{}
	saveUserApp = func(*entity.User) (*entity.User, map[string]string) {
		return nil, map[string]string{
			"email_taken": "email already taken",
		}
	}
	r := gin.Default()
	r.POST("/users", users.SaveUser)
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
////////////////////////////////////////////////////////////////


//GetUsers Test
func TestGetUsers_Success(t *testing.T) {
	//application.UserApp = &fakeUserApp{}
	getUsersApp = func() ([]entity.User, error) {
		//remember we are running sensitive info such as email and password
		return []entity.User{
			 {
				ID:        1,
				FirstName: "victor",
				LastName:  "steven",
			},
			{
				ID:        2,
				FirstName: "mike",
				LastName:  "ken",
			},
		}, nil
	}
	r := gin.Default()
	r.GET("/users", users.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var users []entity.User

	err = json.Unmarshal(rr.Body.Bytes(), &users)

	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, len(users), 2)
}
///////////////////////////////////////////////////////////////



//GetUser Test
func TestGetUser_Success(t *testing.T) {
	//application.UserApp = &fakeUserApp{}
	getUserApp = func(uint64) (*entity.User, error) {
		//remember we are running sensitive info such as email and password
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
	r := gin.Default()
	userId := strconv.Itoa(1)
	r.GET("/users/:user_id", users.GetUser)

	req, err := http.NewRequest(http.MethodGet, "/users/" + userId, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var user *entity.User

	err = json.Unmarshal(rr.Body.Bytes(), &user)

	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, user.FirstName, "victor")
	assert.EqualValues(t, user.LastName, "steven")
}

