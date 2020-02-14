package interfaces

import (
	"bytes"
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/fileupload"
	"food-app/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
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


func TestSaverFood_Success(t *testing.T) {
	application.FoodApp = &fakeFoodApp{}
	token.TokenAuth = &fakeAuth{}
	fileupload.Uploader = &fakeUploader{}
	saveFoodApp = func(*entity.Food) (*entity.Food, map[string]string) {
		//remember we are running sensitive info such as email and password
		return &entity.Food{
			ID:        1,
			UserID:    14,
			Title: "Food title",
			Description:  "Food description",
		}, nil
	}
	fetchAuth = func(ad *token.AccessDetails) (uint64, error){
		return 1, nil
	}
	uploadFile = func(file *multipart.FileHeader) (string, error) {
		return "dbdbfdhbfhbfy3434jhfd.jpg", nil
	}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjYzZjJjNGQyLWFlN2UtNDNmZS05MmNlLWM1Y2VkODgwZTVmOCIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU4MTY5MjA0NCwidXNlcl9pZCI6MTR9.PHLm7XfQa0gvBoSxSEY73S0cwsBx6bocBhBisGYErXg"
	tokenString := fmt.Sprintf("Bearer %v", token)

	file, w := createMultipartFormData(t, "food_image", "amala.jpg")

	data := url.Values{}

	//data := url.
	//data.Add("title", "Food title")
	//data.Add("description", "Food description")
	//data.Add("food_image", file.String())
	r := gin.Default()
	r.POST("/food", SaveFood)

	//req, err := http.NewRequest(http.MethodPost, "/food", strings.NewReader(data.Encode()))
	req, err := http.NewRequest(http.MethodPost, "/food", &file)

	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	//req.Form = data

	req.Header.Set("Authorization", tokenString)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	//fmt.Println("The request: ", req.Header)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	//food := &entity.Food{}

	//err = json.Unmarshal(rr.Body.Bytes(), &food)
	//fmt.Println("error: ", err)
	fmt.Println("response: ", string(rr.Body.Bytes()))
	//fmt.Println("food: ", food)

	//assert.Equal(t, rr.Code, 201)
	//assert.EqualValues(t, user.FirstName, "victor")
	//assert.EqualValues(t, user.LastName, "steven")
}


