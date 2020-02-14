package infrastructure

import (
	"errors"
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"strings"
)

type repositoryFoodCRUD struct {
	db *gorm.DB
}

func NewRepositoryFoodCRUD(db *gorm.DB) repository.FoodRepository {
	return &repositoryFoodCRUD{db}
}

func (r *repositoryFoodCRUD) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&food).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return food, nil
}

func (r *repositoryFoodCRUD) GetFood(id uint64) (*entity.Food, error) {
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

func (r *repositoryFoodCRUD) GetAllFood() ([]entity.Food, error) {
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

func (r *repositoryFoodCRUD) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&food).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["email_taken"] = "email already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return food, nil
}

func (r *repositoryFoodCRUD) DeleteFood(id uint64) error {
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
