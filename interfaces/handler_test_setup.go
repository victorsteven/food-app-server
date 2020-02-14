package interfaces

import "food-app/domain/entity"

var (
	saveUserApp func(*entity.User) (*entity.User, map[string]string)
	getUsersApp func() ([]entity.User, error)
	getUserApp func(uint64) (*entity.User, error)
	getUserEmailPasswordApp func(*entity.User) (*entity.User, map[string]string)

	//Food
	saveFoodApp func(*entity.Food) (*entity.Food, map[string]string)
	getFoodApp func(uint64) (*entity.Food, error)
	updateFoodApp func(*entity.Food) (*entity.Food, map[string]string)

	//getUsersApp func() ([]entity.User, error)
	//getUserApp func(uint64) (*entity.User, error)
	//getUserEmailPasswordApp func(*entity.User) (*entity.User, map[string]string)
)

type fakeUserApp struct {}
type fakeFoodApp struct {}


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
	panic("implement me")
}

func (f *fakeFoodApp) GetFood(foodID uint64) (*entity.Food, error) {
	return getFoodApp(foodID)
}

func (f *fakeFoodApp) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	return updateFoodApp(food)
}

func (f *fakeFoodApp) DeleteFood(uint64) error {
	panic("implement me")
}