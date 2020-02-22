package infrastructure

import (
	"errors"
	"fmt"
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	"os"
	"strings"
)

type foodRepository struct {
	db *gorm.DB
}

//NewRepositoryFood is useful when writing test cases, to swap the real database with a test db
func NewFoodRepository(db *gorm.DB) repository.FoodRepository {
	return &foodRepository{db}
}

var FoodRepo repository.FoodRepository = &foodRepository{}

func (r *foodRepository) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
	db := rdbms.NewDB()
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
	food.FoodImage = os.Getenv("DO_SPACES_URL") + food.FoodImage

	err := db.Debug().Create(&food).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "food title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return food, nil
}

func (r *foodRepository) GetFood(id uint64) (*entity.Food, error) {
	db := rdbms.NewDB()
	var food entity.Food
	err := db.Debug().Where("id = ?", id).Take(&food).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		fmt.Println("error 2: ", err)
		return nil, errors.New("food not found")
	}
	return &food, nil
}

func (r *foodRepository) GetAllFood() ([]entity.Food, error) {
	db := rdbms.NewDB()
	var foods []entity.Food
	err := db.Debug().Limit(100).Order("created_at desc").Find(&foods).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return foods, nil
}

func (r *foodRepository) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
	db := rdbms.NewDB()
	dbErr := map[string]string{}
	err := db.Debug().Save(&food).Error
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

func (r *foodRepository) DeleteFood(id uint64) error {
	db := rdbms.NewDB()
	var food entity.Food
	err := db.Debug().Where("id = ?", id).Delete(&food).Error
	if err != nil {
		return  errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return  errors.New("food not found")
	}
	return  nil
}
