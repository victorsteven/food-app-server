package persistence

import (
    "errors"
    "os"
    "strings"

    "github.com/jinzhu/gorm"

    "food-app/domain/entity"
    "food-app/domain/repository"
)

type FoodRepo struct {
    db *gorm.DB
}

func NewFoodRepository(db *gorm.DB) *FoodRepo {
    return &FoodRepo{db}
}

// FoodRepo implements the repository.FoodRepository interface
var _ repository.FoodRepository = &FoodRepo{}

func (r *FoodRepo) SaveFood(food *entity.Food) (*entity.Food, map[string]string) {
    dbErr := map[string]string{}
    // The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
    food.FoodImage = os.Getenv("DO_SPACES_URL") + food.FoodImage

    err := r.db.Debug().Create(&food).Error
    if err != nil {
        // since our title is unique
        if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
            dbErr["unique_title"] = "food title already taken"
            return nil, dbErr
        }
        // any other db error
        dbErr["db_error"] = "database error"
        return nil, dbErr
    }
    return food, nil
}

func (r *FoodRepo) GetFood(id uint64) (*entity.Food, error) {
    var food entity.Food
    err := r.db.Debug().Where("id = ?", id).Take(&food).Error
    if gorm.IsRecordNotFoundError(err) {
        return nil, errors.New("food not found")
    }
    if err != nil {
        return nil, errors.New("database error, please try again")
    }
    return &food, nil
}

func (r *FoodRepo) GetAllFood() ([]entity.Food, error) {
    var foods []entity.Food
    err := r.db.Debug().Limit(100).Order("created_at desc").Find(&foods).Error
    if gorm.IsRecordNotFoundError(err) {
        return nil, errors.New("user not found")
    }
    if err != nil {
        return nil, err
    }
    return foods, nil
}

func (r *FoodRepo) UpdateFood(food *entity.Food) (*entity.Food, map[string]string) {
    dbErr := map[string]string{}
    err := r.db.Debug().Save(&food).Error
    if err != nil {
        // since our title is unique
        if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
            dbErr["unique_title"] = "title already taken"
            return nil, dbErr
        }
        // any other db error
        dbErr["db_error"] = "database error"
        return nil, dbErr
    }
    return food, nil
}

func (r *FoodRepo) DeleteFood(id uint64) error {
    var food entity.Food
    err := r.db.Debug().Where("id = ?", id).Delete(&food).Error
    if err != nil {
        return errors.New("database error, please try again")
    }
    return nil
}
