package requests

type Register struct {
	Email          string `json:"email"`
	Fullname       string `json:"fullname"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
	Role           string `json:"role"`
	CurrentStage   string `json:"current_stage"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SetRole struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}
