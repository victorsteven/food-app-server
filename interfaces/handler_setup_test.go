package interfaces

import (
	"food-app/domain/entity"
	"food-app/utils/auth"
	"mime/multipart"
	"net/http"
)

//THIS FILE CONTAINS MOCK DATA FROM THE APPLICATION LAYER, THE AUTHENTICATION, ETC. WE DO THIS TO ACHIEVE UNIT TESTING FOR THE CONTROLLER FUNCTIONS

//func TestMain(m *testing.M) {
//	if err := godotenv.Load(os.ExpandEnv("./../.env")); err != nil {
//		log.Println("no env gotten")
//	}
//	os.Exit(m.Run())
//}

var (
	//User
	saveUserApp func(*entity.User) (*entity.User, map[string]string)
	getUsersApp func() ([]entity.User, error)
	getUserApp func(uint64) (*entity.User, error)
	getUserEmailPasswordApp func(*entity.User) (*entity.User, map[string]string)

	//Food
	saveFoodApp func(*entity.Food) (*entity.Food, map[string]string)
	getFoodApp func(uint64) (*entity.Food, error)
	updateFoodApp func(*entity.Food) (*entity.Food, map[string]string)
	uploadFile func(file *multipart.FileHeader) (string, error)
	getAllFoodApp func() ([]entity.Food, error)
	deleteFoodApp func(uint64) error

	tokenMetadata func(*http.Request) (*auth.AccessDetails, error)
	fetchAuth func(uuid string) (uint64, error)
	createAuth func(uint64, *auth.TokenDetails) error

	deleteTokens func(*auth.AccessDetails) error
	deleteRefresh func(string) error
	createToken func(userid uint64) (*auth.TokenDetails, error)

	signin func(*entity.User) (map[string]interface{}, map[string]string)

)

type fakeUserApp struct {}
type fakeFoodApp struct {}
type fakeAuth struct {}
type fakeToken struct {}
type fakeUploader struct {}
type fakeSignin struct {}


//Defining the fake methods.
func (fa *fakeUserApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return getUserEmailPasswordApp(user)
}
func (fa *fakeUserApp) GetUsers() ([]entity.User, error) {
	return getUsersApp()
}
func (fa *fakeUserApp) GetUser(userId uint64) (*entity.User, error) {
	return getUserApp(userId)
}
func (fa *fakeUserApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return saveUserApp(user)
}


func (f *fakeFoodApp) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	return saveFoodApp(food)
}
func (f *fakeFoodApp) GetAllFood() ([]entity.Food, error) {
	return getAllFoodApp()
}
func (f *fakeFoodApp) GetFood(foodID uint64) (*entity.Food, error) {
	return getFoodApp(foodID)
}
func (f *fakeFoodApp) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	return updateFoodApp(food)
}
func (f *fakeFoodApp) DeleteFood(foodId uint64) error {
	return deleteFoodApp(foodId)
}


func (f *fakeAuth) DeleteRefresh(refreshUuid string) error {
	return deleteRefresh(refreshUuid)
}
func (f *fakeAuth) DeleteTokens(authD *auth.AccessDetails) error {
	return deleteTokens(authD)
}
func (f *fakeAuth) FetchAuth(uuid string) (uint64, error) {
	return fetchAuth(uuid)
}
func (f *fakeAuth) CreateAuth(userId uint64, authD *auth.TokenDetails) error {
	return createAuth(userId, authD)
}


func (f *fakeToken) CreateToken(userid uint64) (*auth.TokenDetails, error) {
	return createToken(userid)
}
func (f *fakeToken) ExtractTokenMetadata(r *http.Request) (*auth.AccessDetails, error) {
	return tokenMetadata(r)
}


func (f *fakeUploader) UploadFile(newname *multipart.FileHeader) (string, error) {
	return uploadFile(newname)
}

func (s *fakeSignin) SignIn(user *entity.User) (map[string]interface{}, map[string]string) {
	return signin(user)
}