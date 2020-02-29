package infrastructure

import (
	"fmt"
	"food-app/domain/entity"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func DBConn() (*gorm.DB, error) {
	if _, err := os.Stat("./../.env"); !os.IsNotExist(err) {
		var err error
		err = godotenv.Load(os.ExpandEnv("./../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
		return LocalDatabase()
	}
	return CIBuild()
}

//Circle CI DB
func CIBuild() (*gorm.DB, error) {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "127.0.0.1", "5432", "steven", "food-app-test", "password")
	conn, err := gorm.Open("postgres", DBURL)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
	return conn, nil
}

//Local DB
func LocalDatabase() (*gorm.DB, error) {
	dbdriver := os.Getenv("TEST_DB_DRIVER")
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	user := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_NAME")
	port := os.Getenv("TEST_DB_PORT")

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)
	conn, err := gorm.Open(dbdriver, DBURL)
	if err != nil {
		return nil, err
	} else {
		log.Println("CONNECTED TO: ", dbdriver)
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
		ID:          1,
		Title:       "food title",
		Description: "food desc",
		UserID:      1,
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
			ID:          1,
			Title:       "first food",
			Description: "first desc",
			UserID:      1,
		},
		{
			ID:          2,
			Title:       "second food",
			Description: "second desc",
			UserID:      1,
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
