package models

import (
	"time"
)

type ExpUserSession struct {
	Id        string    `gorm:"default:gen_random_uuid()"`
	ExpUserId string    `gorm:"column:exp_user_id"`
	LastLogin time.Time `gorm:"column:last_login"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
