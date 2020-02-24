package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

type FoodImpl struct {
	infrastructure repository.FoodRepository
}

func NewFoodImpl(inf repository.FoodRepository) *FoodImpl {
	return &FoodImpl{ infrastructure: inf }
}

var _ FoodAppInterface = &FoodImpl{}

type FoodAppInterface interface {
	SaveFood(*entity.Food) (*entity.Food, map[string]string)
	GetAllFood() ([]entity.Food, error)
	GetFood(uint64) (*entity.Food, error)
	UpdateFood(*entity.Food) (*entity.Food, map[string]string)
	DeleteFood(uint64) error
}

func (u *FoodImpl) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	return u.infrastructure.SaveFood(food)
}

func (u *FoodImpl) GetAllFood() ([]entity.Food, error) {
	return u.infrastructure.GetAllFood()
}

func (u *FoodImpl) GetFood(foodId uint64) (*entity.Food, error) {
	return u.infrastructure.GetFood(foodId)
}

func (u *FoodImpl) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	return u.infrastructure.UpdateFood(food)
}

func (u *FoodImpl) DeleteFood(foodId uint64) error {
	return u.infrastructure.DeleteFood(foodId)
}