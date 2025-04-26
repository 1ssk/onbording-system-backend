package responces

type UserResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Fullname     string `json:"fullname"`
	Role         string `json:"role"`
	CurrentStage string `json:"current_stage"`
	Points       int    `json:"points"`
}
