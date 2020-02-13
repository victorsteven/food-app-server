package infrastructure

import (
	"errors"
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type repositoryFood struct {
	db *gorm.DB
}

func NewRepositoryFood(db *gorm.DB) repository.FoodRepository {
	return &repositoryFood{db}
}

func (r *repositoryFood) SaveFood(food *entity.Food) (*entity.Food, error) {
	err := r.db.Debug().Create(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (r *repositoryFood) GetFood(id uint64) (*entity.Food, error) {
	var food entity.Food
	err := r.db.Debug().Where("id = ?", id).Take(&food).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		fmt.Println("error 2: ", err)
		return nil, errors.New("food not found")
	}
	return &food, nil
}

func (r *repositoryFood) GetAllFood() ([]entity.Food, error) {
	var foods []entity.Food
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&foods).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return foods, nil
}

func (r *repositoryFood) UpdateFood(food *entity.Food) (*entity.Food, error) {
	err := r.db.Debug().Save(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (r *repositoryFood) DeleteFood(id uint64) error {
	var food entity.Food
	err := r.db.Debug().Where("id = ?", id).Delete(&food).Error
	if err != nil {
		return  errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return  errors.New("food not found")
	}
	return  nil
}
