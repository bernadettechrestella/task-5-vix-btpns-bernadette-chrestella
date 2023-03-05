package app

import (
	"task5-btpns-stella/models"
)

type Photo struct {
	Id      string
	UserId  string
	Title   string
	Caption string
}

type EditPhoto struct {
	Id      string `json:"id"`
	UserId  string `json:"user_id"`
	Caption string `json:"caption"`
}

type DeletePhoto struct {
	Id     string `json:"id"`
	Userid string `json:"user_id"`
}

type ResponseListPhotos struct {
	Code    string                  `json:"code"`
	Message string                  `json:"message"`
	Data    []models.ExpUserProfile `json:"data"`
}

type GetListPhotos struct {
	UserId   string `json:"user_id"`
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
}
