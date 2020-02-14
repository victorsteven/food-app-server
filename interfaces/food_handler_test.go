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


func TestSaverFood_Success(t *testing.T) {
	application.FoodApp = &fakeFoodApp{} //make it possible to change real method with fake
	token.TokenAuth = &fakeAuth{}
	fileupload.Uploader = &fakeUploader{}

	//Mocking The Food return from db
	saveFoodApp = func(*entity.Food) (*entity.Food, error) {
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


