package application

import (
	"food-app/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS FAILURES

type fakeFoodRepo struct {}

var (
	saveFoodRepo func(*entity.Food) (*entity.Food, map[string]string)
	getFoodRepo func(uint64) (*entity.Food, error)
	getAllFoodRepo func() ([]entity.Food, error)
	updateFoodRepo func(*entity.Food) (*entity.Food, map[string]string)
	deleteFoodRepo func(uint64) error
)

func (f *fakeFoodRepo) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	return saveFoodRepo(food)
}
func (f *fakeFoodRepo) GetFood(foodId uint64) (*entity.Food, error) {
	return getFoodRepo(foodId)
}
func (f *fakeFoodRepo) GetAllFood() ([]entity.Food, error) {
	return getAllFoodRepo()
}
func (f *fakeFoodRepo) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	return updateFoodRepo(food)
}
func (f *fakeFoodRepo) DeleteFood(foodId uint64) error {
	return deleteFoodRepo(foodId)
}

//var fakeFood repository.FoodRepository = &fakeFoodRepo{} //this is where the real implementation is swap with our fake implementation
var foodApp FoodAppInterface = &fakeFoodRepo{} //this is where the real implementation is swap with our fake implementation


func TestSaveFood_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	saveFoodRepo = func(user *entity.Food) (*entity.Food,  map[string]string) {
		return &entity.Food{
			ID:        1,
			Title: "food title",
			Description:  "food description",
			UserID:     1,
		}, nil
	}
	food := &entity.Food{
		ID:        1,
		Title: "food title",
		Description:  "food description",
		UserID:     1,
	}
	f, err := foodApp.SaveFood(food)
	assert.Nil(t, err)
	assert.EqualValues(t, f.Title, "food title")
	assert.EqualValues(t, f.Description, "food description")
	assert.EqualValues(t, f.UserID, 1)
}

func TestGetFood_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getFoodRepo = func(foodId uint64) (*entity.Food,  error) {
		return &entity.Food{
			ID:        1,
			Title: "food title",
			Description:  "food description",
			UserID:     1,
		}, nil
	}
	foodId := uint64(1)
	f, err := foodApp.GetFood(foodId)
	assert.Nil(t, err)
	assert.EqualValues(t, f.Title, "food title")
	assert.EqualValues(t, f.Description, "food description")
	assert.EqualValues(t, f.UserID, 1)
}

func TestAllFood_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	getAllFoodRepo = func() ([]entity.Food,  error) {
		return []entity.Food{
			{
				ID:        1,
				Title: "food title first",
				Description:  "food description first",
				UserID:     1,
			},
			{
				ID:        2,
				Title: "food title second",
				Description:  "food description second",
				UserID:     1,
			},

		}, nil
	}
	f, err := foodApp.GetAllFood()
	assert.Nil(t, err)
	assert.EqualValues(t, len(f), 2)
}

func TestUpdateFood_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	updateFoodRepo = func(user *entity.Food) (*entity.Food,  map[string]string) {
		return &entity.Food{
			ID:        1,
			Title: "food title update",
			Description:  "food description update",
			UserID:     1,
		}, nil
	}
	food := &entity.Food{
		ID:        1,
		Title: "food title update",
		Description:  "food description update",
		UserID:     1,
	}
	f, err := foodApp.UpdateFood(food)
	assert.Nil(t, err)
	assert.EqualValues(t, f.Title, "food title update")
	assert.EqualValues(t, f.Description, "food description update")
	assert.EqualValues(t, f.UserID, 1)
}

func TestDeleteFood_Success(t *testing.T) {
	//Mock the response coming from the infrastructure
	deleteFoodRepo = func(foodId uint64) error {
		return nil
	}
	foodId := uint64(1)
	err := foodApp.DeleteFood(foodId)
	assert.Nil(t, err)
}
