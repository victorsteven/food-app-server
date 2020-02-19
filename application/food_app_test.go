package application

import (
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
	"github.com/stretchr/testify/assert"
	"testing"
)

//IF YOU HAVE TIME, YOU CAN TEST ALL THE METHODS

type fakeFoodRepo struct {}

var (
	saveFoodRepo func(*entity.Food) (*entity.Food, map[string]string)
)

func (f fakeFoodRepo) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	return saveFoodRepo(food)
}

func (f fakeFoodRepo) GetFood(uint64) (*entity.Food, error) {
	panic("implement me")
}

func (f fakeFoodRepo) GetAllFood() ([]entity.Food, error) {
	panic("implement me")
}

func (f fakeFoodRepo) UpdateFood(*entity.Food) (*entity.Food, map[string]string) {
	panic("implement me")
}

func (f fakeFoodRepo) DeleteFood(uint64) error {
	panic("implement me")
}

func TestSaveFood_Success(t *testing.T) {
	infrastructure.FoodRepo = &fakeFoodRepo{}
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
	f, err := FoodApp.SaveFood(food)
	assert.Nil(t, err)
	assert.EqualValues(t, f.Title, "food title")
	assert.EqualValues(t, f.Description, "food description")
	assert.EqualValues(t, f.UserID, 1)
}
