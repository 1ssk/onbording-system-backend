package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `gorm:"unique"`
	Email      string `gorm:"unique;not null"`
	Password   string
	Role       string `gorm:"default:user"`
}
