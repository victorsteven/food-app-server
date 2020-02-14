package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/fileupload"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load(os.ExpandEnv("./../.env")); err != nil {
		log.Println("no env gotten")
	}
	os.Exit(m.Run())
}

var (
	fetchAuth func(ad *token.AccessDetails) (uint64, error)
	uploadFile func(file *multipart.FileHeader) (string, error)
)
type fakeAuth struct {}
type fakeUploader struct {}

func (f *fakeUploader) UploadFile(newname *multipart.FileHeader) (string, error) {
	return uploadFile(newname)
}

func (f *fakeAuth) CreateAuth(uint64, *token.TokenDetails) error {
	panic("implement me")
}

func (f *fakeAuth) NewRedisClient(host, port, password string) (*redis.Client, error) {
	panic("implement me")
}

func (f *fakeAuth) FetchAuth(ad *token.AccessDetails) (uint64, error){
	return fetchAuth(ad)
}

func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
	var b bytes.Buffer
	var err error
	w := multipart.NewWriter(&b)
	var fw io.Writer
	file := mustOpen(fileName)
	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
		t.Errorf("Error creating writer: %v", err)
	}
	if _, err = io.Copy(fw, file); err != nil {
		t.Errorf("Error with io.Copy: %v", err)
	}
	w.Close()
	return b, w
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		pwd, _ := os.Getwd()
		fmt.Println("PWD: ", pwd)
		panic(err)
	}
	return r
}

func Test_SaveFood_Invalid_Data(t *testing.T) {
	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			inputJSON:  `{"title": "", "description": "steven","email": "steven@example.com","password": "password"}`,
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
	}
}


func TestSaverFood_Success(t *testing.T) {
	application.FoodApp = &fakeFoodApp{} //make it possible to change real method with fake
	token.TokenAuth = &fakeAuth{}
	fileupload.Uploader = &fakeUploader{}

	//Mocking The Food return from db
	saveFoodApp = func(*entity.Food) (*entity.Food, map[string]string) {
		return &entity.Food{
			ID:        1,
			UserID:    1,
			Title: "Food title",
			Description:  "Food description",
			FoodImage: "dbdbf-dhbfh-bfy34-34jh-fd.jpg",

		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fetchAuth = func(ad *token.AccessDetails) (uint64, error){
		return 1, nil
	}
	//Mocking file upload to DigitalOcean
	uploadFile = func(file *multipart.FileHeader) (string, error) {
		return "dbdbf-dhbfh-bfy34-34jh-fd.jpg", nil //this is fabricated
	}

	f :=  "./../utils/images/amala.jpg" //this is where the image is located
	file, err := os.Open(f)
	if err != nil {
		t.Errorf("Cannot open file: %s\n", err)
	}
	defer file.Close()

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Initialize the file field
	fileWriter, err := multipartWriter.CreateFormFile("food_image", "amala.jpg")
	if err != nil {
		t.Errorf("Cannot write file: %s\n", err)
	}
	//Copy the actual content to the file field's writer, though this is not needed, since files are sent to DigitalOcean
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Errorf("Cannot copy file: %s\n", err)
	}
	//Add the title and the description fields
	fileWriter, err = multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("Cannot write title: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Food title"))
	if err != nil {
		t.Errorf("Cannot write title value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("description")
	if err != nil {
		t.Errorf("Cannot write description: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Food description"))
	if err != nil {
		t.Errorf("Cannot write description value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest(http.MethodPost, "/food", &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.POST("/food", SaveFood)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var food = entity.Food{}
	err = json.Unmarshal(rr.Body.Bytes(), &food)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 201)
	assert.EqualValues(t, food.ID, 1)
	assert.EqualValues(t, food.UserID, 1)
	assert.EqualValues(t, food.Title, "Food title")
	assert.EqualValues(t, food.Description, "Food description")
	assert.EqualValues(t, food.FoodImage, "dbdbf-dhbfh-bfy34-34jh-fd.jpg")
}


func TestUpdateFood_Success_With_File(t *testing.T) {
	application.FoodApp = &fakeFoodApp{} //make it possible to change real method with fake
	token.TokenAuth = &fakeAuth{}
	fileupload.Uploader = &fakeUploader{}

	//Return Food to check for, with our mock
	getFoodApp = func(uint64) (*entity.Food, error) {
		return &entity.Food{
			ID:        1,
			UserID:    1,
			Title: "Food title",
			Description:  "Food description",
			FoodImage: "dbdbf-dhbfh-bfy34-34jh-fd.jpg",
		}, nil
	}
	//Mocking The Food return from db
	updateFoodApp = func(*entity.Food) (*entity.Food, map[string]string) {
		return &entity.Food{
			ID:        1,
			UserID:    1,
			Title: "Food title updated",
			Description:  "Food description updated",
			FoodImage: "dbdbf-dhbfh-bfy34-34jh-fd-updated.jpg",
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fetchAuth = func(ad *token.AccessDetails) (uint64, error){
		return 1, nil
	}
	//Mocking file upload to DigitalOcean
	uploadFile = func(file *multipart.FileHeader) (string, error) {
		return "dbdbf-dhbfh-bfy34-34jh-fd-updated.jpg", nil //this is fabricated
	}

	f :=  "./../utils/images/new_meal.jpeg" //this is where the image is located
	file, err := os.Open(f)
	if err != nil {
		t.Errorf("Cannot open file: %s\n", err)
	}
	defer file.Close()

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Initialize the file field
	fileWriter, err := multipartWriter.CreateFormFile("food_image", "new_meal.jpeg")
	if err != nil {
		t.Errorf("Cannot write file: %s\n", err)
	}
	//Copy the actual content to the file field's writer, though this is not needed, since files are sent to DigitalOcean
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		t.Errorf("Cannot copy file: %s\n", err)
	}
	//Add the title and the description fields
	fileWriter, err = multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("Cannot write title: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Food title updated"))
	if err != nil {
		t.Errorf("Cannot write title value: %s\n", err)
	}
	fileWriter, err = multipartWriter.CreateFormField("description")
	if err != nil {
		t.Errorf("Cannot write description: %s\n", err)
	}
	_, err = fileWriter.Write([]byte("Food description updated"))
	if err != nil {
		t.Errorf("Cannot write description value: %s\n", err)
	}
	//Close the multipart writer so it writes the ending boundary
	multipartWriter.Close()

	//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"

	tokenString := fmt.Sprintf("Bearer %v", token)

	foodID := strconv.Itoa(1)
	req, err := http.NewRequest(http.MethodPut, "/food/"+foodID, &requestBody)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	r := gin.Default()
	r.PUT("/food/:food_id", UpdateFood)
	req.Header.Set("Authorization", tokenString)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType()) //this is important
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var food = entity.Food{}
	err = json.Unmarshal(rr.Body.Bytes(), &food)
	if err != nil {
		t.Errorf("cannot unmarshal response: %v\n", err)
	}
	assert.Equal(t, rr.Code, 200)
	assert.EqualValues(t, food.ID, 1)
	assert.EqualValues(t, food.UserID, 1)
	assert.EqualValues(t, food.Title, "Food title updated")
	assert.EqualValues(t, food.Description, "Food description updated")
	assert.EqualValues(t, food.FoodImage, "dbdbf-dhbfh-bfy34-34jh-fd-updated.jpg")
}