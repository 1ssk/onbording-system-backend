package requests

type CompleteCourse struct {
	Email string `json:"email"`
}

type AddMeeting struct {
	UserID    string `json:"user_id"`
	LeaderID  string `json:"leader_id"`
	Title     string `json:"title"`
	StartTime string `json:"start_time"`
}

type UpdateMeeting struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	LeaderID string `json:"leader_id"`
	Title    string `json:"title"`
}

type DeleteMeeting struct {
	ID string `json:"id"`
}

type AddUser struct {
	Email        string `json:"email"`
	Fullname     string `json:"fullname"`
	Password     string `json:"password"`
	CurrentStage string `json:"current_stage"`
}
