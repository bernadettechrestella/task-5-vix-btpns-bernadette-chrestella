package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	l "task5-btpns-stella/app"
	"task5-btpns-stella/models"
	"time"
)

func UserUpdateController(c *gin.Context) {
	var req l.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Field is Empty")
		return
	}

	_, payload, err := authorized(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if req.Id != "" && req.NewEmail != "" && req.OldEmail != "" && req.Id == payload.UserId {

		user := &models.ExpUser{}

		dbcontext.Table("exp_user").First(user, "id = ?", payload.UserId)

		if user.Email == req.OldEmail {
			user.Email = req.NewEmail
			user.UpdatedAt = time.Now()

			err := dbcontext.Table("exp_user").Save(user).Error
			if err != nil {
				errorResponse(c, http.StatusInternalServerError, err.Error())
				return
			}

			res := l.ApiResponse{
				Message: "Success",
				Code:    "200",
			}

			c.JSON(http.StatusOK, res)

		} else {
			res := l.ApiResponse{
				Message: "Email yang terdaftar sekarang tidak sesuai dengan yang di masukan",
				Code:    "500",
			}

			c.JSON(http.StatusInternalServerError, res)
		}

		//userInfo, err := client.GetUserInfo(ctx, token, "master")
		//if err != nil {
		//	errorResponse(c, http.StatusInternalServerError, err.Error())
		//	return
		//}
		//
		//err = client.SetPassword(ctx, token, req.Id, "master", req.Password, false)
		//if err != nil {
		//	errorResponse(c, http.StatusInternalServerError, err.Error())
		//	return
		//}
	}

}

func UserDeleteController(c *gin.Context) {
	var req l.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Field is Empty")
		return
	}

	_, payload, err := authorized(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if req.Id != "" {
		user := &models.ExpUser{}

		dbcontext.Table("exp_user").First(user, "id = ?", payload.UserId)

		user.Active = false
		user.UpdatedAt = time.Now()

		err := dbcontext.Table("exp_user").Save(user).Error
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		res := l.ApiResponse{
			Message: "Success",
			Code:    "200",
		}

		c.JSON(http.StatusOK, res)
	}
}
