package models

import "time"

type ExpUserProfile struct {
	Id        string    `gorm:"default:gen_random_uuid()"`
	ExpUserId string    `gorm:"column:exp_user_id"`
	Title     string    `gorm:"column:title"`
	Caption   string    `gorm:"column:caption"`
	PhotoUrl  string    `gorm:"column:photo_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Deleted   bool      `gorm:"column:is_deleted"`
}
