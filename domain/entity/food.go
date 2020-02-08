package entity

import (
	"os"
	"time"
)

type Food struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint64     `gorm:"size:100;not null;" json:"user_id"`
	Title       string     `gorm:"size:100;not null;" json:"title"`
	Description string     `gorm:"text;not null;" json:"description"`
	AvatarPath string    `gorm:"size:255;null;" json:"avatar_path"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (f *Food) AfterFind()  {
	if f.AvatarPath != "" {
		f.AvatarPath = os.Getenv("DO_SPACES_URL") + f.AvatarPath
	}
}