package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"food-app/application"
	"food-app/domain/entity"
	"food-app/utils/fileupload"
	"food-app/utils/auth"
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
	tokenMetadata func(*http.Request) (*auth.AccessDetails, error)
	fetchAuth func(uuid string) (uint64, error)
	uploadFile func(file *multipart.FileHeader) (string, error)
)
type fakeAuth struct {}

func (f *fakeAuth) DeleteRefresh(string) (int64, error) {
	panic("implement me")
}

func (f *fakeAuth) DeleteTokens(*auth.AccessDetails) (int64, error) {
	panic("implement me")
}

type fakeToken struct {}

func (f *fakeToken) CreateToken(userid uint64) (*auth.TokenDetails, error) {
	panic("implement me")
}

type fakeUploader struct {}

func (f *fakeAuth) FetchAuth(uuid string) (uint64, error) {
	return fetchAuth(uuid)
}
func (f *fakeToken) ExtractTokenMetadata(r *http.Request) (*auth.AccessDetails, error) {
	return tokenMetadata(r)
}

func (f *fakeUploader) UploadFile(newname *multipart.FileHeader) (string, error) {
	return uploadFile(newname)
}
func (f *fakeAuth) CreateAuth(uint64, *auth.TokenDetails) error {
	panic("implement me")
}
func (f *fakeAuth) NewRedisClient(host, port, password string) (*redis.Client, error) {
	panic("implement me")
}

//func (f *fakeToken) ExtractAndValidateMetadata(r *http.Request) (uint64, error){
//	return tokenMetadata(r)
//}

//func createMultipartFormData(t *testing.T, fieldName, fileName string) (bytes.Buffer, *multipart.Writer) {
//	var b bytes.Buffer
//	var err error
//	w := multipart.NewWriter(&b)
//	var fw io.Writer
//	file := mustOpen(fileName)
//	if fw, err = w.CreateFormFile(fieldName, file.Name()); err != nil {
//		t.Errorf("Error creating writer: %v", err)
//	}
//	if _, err = io.Copy(fw, file); err != nil {
//		t.Errorf("Error with io.Copy: %v", err)
//	}
//	w.Close()
//	return b, w
//}

//func mustOpen(f string) *os.File {
//	r, err := os.Open(f)
//	if err != nil {
//		pwd, _ := os.Getwd()
//		fmt.Println("PWD: ", pwd)
//		panic(err)
//	}
//	return r
//}

func Test_SaveFood_Invalid_Data(t *testing.T) {

	auth.Token = &fakeToken{}
	auth.Auth = &fakeAuth{}

	//Mock extracting metadata
	tokenMetadata = func(r *http.Request) (*auth.AccessDetails, error){
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fetchAuth =  func(uuid string) (uint64, error){
		return 1, nil
	}

	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			//when the title is empty
			inputJSON:  `{"title": "", "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//the description is empty
			inputJSON:  `{"title": "the title", "description": ""}`,
			statusCode: 422,
		},
		{
			//both the title and the description are empty
			inputJSON:  `{"title": "", "description": ""}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": 12344, "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": "hello title", "description": 3242342}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"
		tokenString := fmt.Sprintf("Bearer %v", token)

		r := gin.Default()
		r.POST("/food", SaveFood)
		req, err := http.NewRequest(http.MethodPost, "/food", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req.Header.Set("Authorization", tokenString)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		validationErr := make(map[string]string)

		err = json.Unmarshal(rr.Body.Bytes(), &validationErr)
		if err != nil {
			t.Errorf("error unmarshalling error %s\n", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if validationErr["title_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
		}
		if validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["title_required"] != "" && validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}


func TestSaverFood_Success(t *testing.T) {
	application.FoodApp = &fakeFoodApp{} //make it possible to change real method with fake
	auth.Token = &fakeToken{}
	auth.Auth = &fakeAuth{}
	application.UserApp = &fakeUserApp{}
	fileupload.Uploader = &fakeUploader{}

	//Mock extracting metadata
	tokenMetadata = func(r *http.Request) (*auth.AccessDetails, error){
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fetchAuth =  func(uuid string) (uint64, error){
		return 1, nil
	}
	getUserApp = func(uint64) (*entity.User, error) {
		//remember we are running sensitive info such as email and password
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
	//Mocking file upload to DigitalOcean
	uploadFile = func(file *multipart.FileHeader) (string, error) {
		return "dbdbf-dhbfh-bfy34-34jh-fd.jpg", nil //this is fabricated
	}
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
	auth.Token = &fakeToken{}
	auth.Auth = &fakeAuth{}
	application.UserApp = &fakeUserApp{}
	fileupload.Uploader = &fakeUploader{}

	//Mock extracting metadata
	tokenMetadata = func(r *http.Request) (*auth.AccessDetails, error){
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fetchAuth =  func(uuid string) (uint64, error){
		return 1, nil
	}
	getUserApp = func(uint64) (*entity.User, error) {
		//remember we are running sensitive info such as email and password
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
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


//This is where file is not updated. A user can choose not to update file, in that case, the old file will still be used
func TestUpdateFood_Success_Without_File(t *testing.T) {
	application.FoodApp = &fakeFoodApp{} //make it possible to change real method with fake
	auth.Token = &fakeToken{}
	auth.Auth = &fakeAuth{}
	application.UserApp = &fakeUserApp{}
	fileupload.Uploader = &fakeUploader{}

	//Mock extracting metadata
	tokenMetadata = func(r *http.Request) (*auth.AccessDetails, error){
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fetchAuth =  func(uuid string) (uint64, error){
		return 1, nil
	}
	getUserApp = func(uint64) (*entity.User, error) {
		//remember we are running sensitive info such as email and password
		return &entity.User{
			ID:        1,
			FirstName: "victor",
			LastName:  "steven",
		}, nil
	}
	//Return Food to check for, with our mock
	getFoodApp = func(uint64) (*entity.Food, error) {
		return &entity.Food{
			ID:        1,
			UserID:    1,
			Title: "Food title",
			Description:  "Food description",
			FoodImage: "dbdbf-dhbfh-bfy34-34jh-fd-old-file.jpg",
		}, nil
	}
	//Mocking The Food return from db
	updateFoodApp = func(*entity.Food) (*entity.Food, map[string]string) {
		return &entity.Food{
			ID:        1,
			UserID:    1,
			Title: "Food title updated",
			Description:  "Food description updated",
			FoodImage: "dbdbf-dhbfh-bfy34-34jh-fd-old-file.jpg",
		}, nil
	}

	//Mocking file upload to DigitalOcean
	uploadFile = func(file *multipart.FileHeader) (string, error) {
		return "dbdbf-dhbfh-bfy34-34jh-fd-old-file.jpg", nil //this is fabricated
	}

	//Create a buffer to store our request body as bytes
	var requestBody bytes.Buffer

	//Create a multipart writer
	multipartWriter := multipart.NewWriter(&requestBody)

	//Add the title and the description fields
	fileWriter, err := multipartWriter.CreateFormField("title")
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
	assert.EqualValues(t, food.FoodImage, "dbdbf-dhbfh-bfy34-34jh-fd-old-file.jpg")
}


func Test_UdpateFood_Invalid_Data(t *testing.T) {
	auth.Token = &fakeToken{}
	auth.Auth = &fakeAuth{}
	//Mocking the fetching of token metadata from redis

	//Mock extracting metadata
	tokenMetadata = func(r *http.Request) (*auth.AccessDetails, error){
		return &auth.AccessDetails{
			TokenUuid: "0237817a-1546-4ca3-96a4-17621c237f6b",
			UserId:    1,
		}, nil
	}
	//Mocking the fetching of token metadata from redis
	fetchAuth =  func(uuid string) (uint64, error){
		return 1, nil
	}

	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			//when the title is empty
			inputJSON:  `{"title": "", "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//the description is empty
			inputJSON:  `{"title": "the title", "description": ""}`,
			statusCode: 422,
		},
		{
			//both the title and the description are empty
			inputJSON:  `{"title": "", "description": ""}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": 12344, "description": "the desc"}`,
			statusCode: 422,
		},
		{
			//When invalid data is passed, e.g, instead of an integer, a string is passed
			inputJSON:  `{"title": "hello sir", "description": 3242342}`,
			statusCode: 422,
		},
	}

	for _, v := range samples {

		//use a valid token that has not expired. This token was created to live forever, just for test purposes with the user id of 1. This is so that it can always be used to run tests
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjgyYTM3YWE5LTI4MGMtNDQ2OC04M2RmLTZiOGYyMDIzODdkMyIsImF1dGhvcml6ZWQiOnRydWUsInVzZXJfaWQiOjF9.ESelxq-UHormgXUwRNe4_Elz2i__9EKwCXPsNCyKV5o"
		tokenString := fmt.Sprintf("Bearer %v", token)

		r := gin.Default()
		r.POST("/food", SaveFood)
		req, err := http.NewRequest(http.MethodPost, "/food", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req.Header.Set("Authorization", tokenString)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		validationErr := make(map[string]string)

		err = json.Unmarshal(rr.Body.Bytes(), &validationErr)
		if err != nil {
			t.Errorf("error unmarshalling error %s\n", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)


		if validationErr["title_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
		}
		if validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["title_required"] != "" && validationErr["description_required"] != "" {
			assert.Equal(t, validationErr["title_required"], "title is required")
			assert.Equal(t, validationErr["description_required"], "description is required")
		}
		if validationErr["invalid_json"] != "" {
			assert.Equal(t, validationErr["invalid_json"], "invalid json")
		}
	}
}