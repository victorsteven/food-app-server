package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type FoodApp struct {
	fr repository.FoodRepository
}


var _ FoodAppInterface = &FoodApp{}

type FoodAppInterface interface {
	SaveFood(*entity.Food) (*entity.Food, map[string]string)
	GetAllFood() ([]entity.Food, error)
	GetFood(uint64) (*entity.Food, error)
	UpdateFood(*entity.Food) (*entity.Food, map[string]string)
	DeleteFood(uint64) error
}

func (f *FoodApp) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	return f.fr.SaveFood(food)
}

func (f *FoodApp) GetAllFood() ([]entity.Food, error) {
	return f.fr.GetAllFood()
}

func (f *FoodApp) GetFood(foodId uint64) (*entity.Food, error) {
	return f.fr.GetFood(foodId)
}

func (f *FoodApp) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	return f.fr.UpdateFood(food)
}

func (f *FoodApp) DeleteFood(foodId uint64) error {
	return f.fr.DeleteFood(foodId)
}
