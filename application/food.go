package application

import (
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"food-app/domain/infrastructure"
)

type FoodImpl struct {
}

var FoodApp FoodAppInterface = &FoodImpl{}

type FoodAppInterface interface {
	SaveFood(*entity.Food) (*entity.Food, map[string]string)
	GetAllFood() ([]entity.Food, error)
	GetFood(uint64) (*entity.Food, error)
	UpdateFood(*entity.Food) (*entity.Food, map[string]string)
	DeleteFood(uint64) error
}

func (u *FoodImpl) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryFoodCRUD(db)
	return conn.SaveFood(food)
}

func (u *FoodImpl) GetAllFood() ([]entity.Food, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryFoodCRUD(db)
	return conn.GetAllFood()
}

func (u *FoodImpl) GetFood(foodId uint64) (*entity.Food, error) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryFoodCRUD(db)
	return conn.GetFood(foodId)
}

func (u *FoodImpl) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryFoodCRUD(db)
	return conn.UpdateFood(food)
}

func (u *FoodImpl) DeleteFood(foodId uint64) error {
	db := rdbms.NewDB()
	conn := infrastructure.NewRepositoryFoodCRUD(db)
	return conn.DeleteFood(foodId)
}