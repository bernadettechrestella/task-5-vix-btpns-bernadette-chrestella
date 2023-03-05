package app

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CountLogin struct {
	Count int `json:"count"`
}
