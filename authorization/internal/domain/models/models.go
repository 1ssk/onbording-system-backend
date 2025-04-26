package models

const (
	RoleUser   = "user"
	RoleLeader = "leader"
	RoleMentor = "mentor"
	RoleAdmin  = "admin"
)

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Fullname     string `json:"fullname"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	CurrentStage string `json:"current_stage"`
	Points       int    `json:"points"`
}
