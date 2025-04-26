package models

import (
	"time"

	"gorm.io/gorm"
)

type Meeting struct {
	gorm.Model `gorm:"unique"`
	UserID    string     `json:"user_id"`
	LeaderID  string     `json:"leader_id"`
	Title     string     `json:"title"`
	StartTime *time.Time `json:"start_time"`
	Status    string     `json:"status"`
	Feedback  string     `json:"feedback"`
}

type User struct {
	gorm.Model   `gorm:"unique"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
}
