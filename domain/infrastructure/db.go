package infrastructure

import (
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"


)

type Repositories struct {
	User repository.UserRepository
	Food    repository.FoodRepository
	db      *gorm.DB
}

func NewServices(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Dbdriver, DBURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User: NewUserRepository(db),
		Food: NewFoodRepository(db),
		db:   db,
	}, nil
}

//closes the  database connection
func (r *Repositories) Close() error {
	return r.db.Close()
}

//This migrate all tables
func (r *Repositories) Automigrate() error {
	return r.db.AutoMigrate(&entity.User{}, &entity.Food{}).Error
}
