package responces

type MeatingResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	LeaderID  string `json:"leader_id"`
	Title     string `json:"title"`
	StartTime string `json:"start_time"`
	Status    string `json:"status"`
	Feedback  string `json:"feedback"`
}
