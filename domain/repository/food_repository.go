package repository

import "food-app/domain/entity"

type FoodRepository interface {
	SaveFood(*entity.Food) (*entity.Food, error)
	GetFood(uint64) (*entity.Food, error)
	GetAllFood() ([]entity.Food, error)
	UpdateFood(*entity.Food) (*entity.Food, error)
}
