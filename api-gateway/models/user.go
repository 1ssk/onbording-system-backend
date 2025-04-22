package models

type User struct {
    ID    uint   `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}
