package models

import "time"

type ExpUser struct {
	Id        string    `gorm:"primaryKey;column:id;default:gen_random_uuid()"`
	Username  string    `gorm:"column:username"`
	Email     string    `gorm:"column:email"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Active    bool      `gorm:"column:active"`
}

type SuccessResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
