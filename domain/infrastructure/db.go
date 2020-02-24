package infrastructure

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
	"github.com/jinzhu/gorm"
)

type Services struct {
	User repository.UserRepository
	//User    UserService
	//Image   ImageService
	db      *gorm.DB
}

func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Services{
		User:    NewUserService(db),
		//Gallery: NewGalleryService(db),
		//Image:   NewImageService(),
		db:      db,
	}, nil
}

//closes the  database connection
func (s *Services) Close() error {
	return s.db.Close()
}

//Drops all tables and rebuild them
//func (s *Services) DestructiveReset() error {
//	if err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error; err != nil {
//		return err
//	}
//	return s.Automigrate()
//}

//This migrate all tables
func (s *Services) Automigrate() error {
	return s.db.AutoMigrate(&entity.User{}, &entity.Food{}).Error
}
