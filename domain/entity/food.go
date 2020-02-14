package entity

import (
	"html"
	"os"
	"strings"
	"time"
)

type Food struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserID      uint64     `gorm:"size:100;not null;" json:"user_id"`
	Title       string     `gorm:"size:100;not null;unique" json:"title"`
	Description string     `gorm:"text;not null;" json:"description"`
	FoodImage string    `gorm:"size:255;null;" json:"food_image"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	//User    User      `json:"user"` //a user can create many food. This also means, anytime you are getting food, also get the user
}


func (f *Food) BeforeSave() {
	f.Title = html.EscapeString(strings.TrimSpace(f.Title))
}

func (f *Food) AfterFind()  {
	if f.FoodImage != "" {
		f.FoodImage = os.Getenv("DO_SPACES_URL") + f.FoodImage
	}
}
func (f *Food) AfterSave()  {
	if f.FoodImage != "" {
		f.FoodImage = os.Getenv("DO_SPACES_URL") + f.FoodImage
	}
}
func (f *Food) AfterUpdate()  {
	if f.FoodImage != "" {
		f.FoodImage = os.Getenv("DO_SPACES_URL") + f.FoodImage
	}
}

func (f *Food) Prepare() {
	f.Title = html.EscapeString(strings.TrimSpace(f.Title))
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()
}

func (f *Food) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if f.Title == "" {
			errorMessages["title_required"] = "title required"
		}
		if f.Description == "" {
			errorMessages["desc_required"] = "description required"
		}
	default:
		if f.Title == "" {
			errorMessages["title_required"] = "title required"
		}
		if f.Description == "" {
			errorMessages["desc_required"] = "description required"
		}
	}
	return errorMessages
}