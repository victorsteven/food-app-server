package infrastructure

import (
	"food-app/database/rdbms"
	"food-app/domain/entity"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(os.ExpandEnv("./../../.env")); err != nil {
		log.Println("no env gotten")
	}
}
func Database() (*gorm.DB, error) {
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	conn, err := rdbms.NewDBConnection(dbdriver, user, password, port, host, dbname)
	if err != nil {
		return nil, err
	}
	err = conn.DropTableIfExists(&entity.User{}, &entity.Food{}).Error
	if err != nil {
		return nil, err
	}
	err = conn.Debug().AutoMigrate(
		entity.User{},
		entity.Food{},
	).Error
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func seedUser(db *gorm.DB) (*entity.User, error) {
	user := &entity.User{
		ID:        1,
		FirstName: "vic",
		LastName:  "stev",
		Email:     "steven@example.com",
		Password:  "password",
		DeletedAt: nil,
	}
	err := db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func seedUsers(db *gorm.DB) ([]entity.User, error) {
	users := []entity.User{
		{
			ID:        1,
			FirstName: "vic",
			LastName:  "stev",
			Email:     "steven@example.com",
			Password:  "password",
			DeletedAt: nil,
		},
		{
			ID:        2,
			FirstName: "kobe",
			LastName:  "bryant",
			Email:     "kobe@example.com",
			Password:  "password",
			DeletedAt: nil,
		},
	}
	for _, v := range users {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}

func seedFood(db *gorm.DB) (*entity.Food, error) {
	food := &entity.Food{
		ID:        1,
		Title: "food title",
		Description:  "food desc",
		UserID:     1,
	}
	err := db.Create(&food).Error
	if err != nil {
		return nil, err
	}
	return food, nil
}

func seedFoods(db *gorm.DB) ([]entity.Food, error) {
	foods := []entity.Food{
		{
			ID:        1,
			Title: "first food",
			Description:  "first desc",
			UserID:     1,
		},
		{
			ID:        2,
			Title: "second food",
			Description:  "second desc",
			UserID:     1,
		},
	}
	for _, v := range foods {
		err := db.Create(&v).Error
		if err != nil {
			return nil, err
		}
	}
	return foods, nil
}


