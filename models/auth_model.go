package models

import "github.com/jinzhu/gorm"

type LoginInput struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}
