package app

type UserUpdate struct {
	Id       string `json:"id"`
	OldEmail string `json:"old_email"`
	NewEmail string `json:"new_email"`
}
