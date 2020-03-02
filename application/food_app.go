package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type foodApp struct {
	fr repository.FoodRepository
}


var _ FoodAppInterface = &foodApp{}

type FoodAppInterface interface {
	SaveFood(*entity.Food) (*entity.Food, map[string]string)
	GetAllFood() ([]entity.Food, error)
	GetFood(uint64) (*entity.Food, error)
	UpdateFood(*entity.Food) (*entity.Food, map[string]string)
	DeleteFood(uint64) error
}

func (f *foodApp) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	return f.fr.SaveFood(food)
}

func (f *foodApp) GetAllFood() ([]entity.Food, error) {
	return f.fr.GetAllFood()
}

func (f *foodApp) GetFood(foodId uint64) (*entity.Food, error) {
	return f.fr.GetFood(foodId)
}

func (f *foodApp) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	return f.fr.UpdateFood(food)
}

func (f *foodApp) DeleteFood(foodId uint64) error {
	return f.fr.DeleteFood(foodId)
}
