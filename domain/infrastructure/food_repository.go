package infrastructure

import (
	"errors"
	"food-app/domain/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type repositoryFoodCRUD struct {
	db *gorm.DB
}

func NewRepositoryFoodCRUD(db *gorm.DB) *repositoryFoodCRUD {
	return &repositoryFoodCRUD{db}
}

//var (
//	UserRepo repository.UserRepository = &userRepo{}
//)

//func NewServer(db *gorm.DB) repository.UserRepository {
//	return &Server{DB: db}
//}

func (r *repositoryFoodCRUD) SaveFood(food *entity.Food) (*entity.Food, error) {
	err := r.db.Debug().Create(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil
}

func (r *repositoryFoodCRUD) GetFood(id uint64) (*entity.Food, error) {
	var food entity.Food
	err := r.db.Debug().Where("id = ?", id).Take(&food).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &food, nil
}

func (r *repositoryFoodCRUD) GetAllFood() ([]entity.Food, error) {
	var allfood []entity.Food
	err := r.db.Debug().Find(&allfood).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return allfood, nil
}
