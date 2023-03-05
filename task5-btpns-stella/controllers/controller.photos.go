package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	l "task5-btpns-stella/app"
	"task5-btpns-stella/models"
	"time"
)

func PhotoPostController(c *gin.Context) {
	var req l.Photo

	req.UserId = c.PostForm("user_id")
	req.Title = c.PostForm("title")
	req.Caption = c.PostForm("caption")

	f, err := c.FormFile("photo")
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Status Unauthorized")
		return
	}

	_, payload, err := authorized(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if payload.UserId == req.UserId {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Second*50)

		defer cancel()

		bucket := os.Getenv("FIREBASE_CLOUD_STORAGE")
		bucketName := strings.Split(bucket, "gs://")

		client, err := InitFirebaseStorage()
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		file, err := f.Open()
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		var errs error
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			wc := client.Bucket(bucketName[1]).Object(f.Filename).NewWriter(ctx)
			if _, err = io.Copy(wc, file); err != nil {
				errs = err
			}
			if err := wc.Close(); err != nil {
				errs = err
			}
			wg.Done()
		}()
		wg.Wait()
		if errs != nil {
			errorResponse(c, http.StatusInternalServerError, errs.Error())
			return
		}

		photo := &models.ExpUserProfile{
			ExpUserId: req.UserId,
			PhotoUrl:  bucket + "/" + f.Filename,
			Title:     req.Title,
			Caption:   req.Caption,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Deleted:   false,
		}

		err = dbcontext.Table("exp_user_profile").
			Create(photo).Error
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, errs.Error())
			return
		}

		res := l.ApiResponse{
			Message: "Success",
			Code:    "200",
		}

		c.JSON(http.StatusOK, res)

	} else {
		res := l.ApiResponse{
			Message: "Unauthorized",
			Code:    "401",
		}

		c.JSON(http.StatusUnauthorized, res)
	}

}

func PhotoGetController(c *gin.Context) {
	var req l.GetListPhotos
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Field is Empty")
		return
	}

	_, payload, err := authorized(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if payload.UserId == req.UserId {
		var result []models.ExpUserProfile

		err := dbcontext.Scopes(Paginate(req)).
			Table("exp_user_profile").
			Where("exp_user_id = ?", req.UserId).
			Where("is_deleted = ?", false).
			Find(&result).Error
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		res := l.ResponseListPhotos{
			Code:    "200",
			Message: "Success",
			Data:    result,
		}

		c.JSON(http.StatusOK, res)

	} else {
		res := l.ApiResponse{
			Message: "Unauthorized",
			Code:    "401",
		}

		c.JSON(http.StatusUnauthorized, res)
	}
}

func PhotoEditController(c *gin.Context) {
	var req l.EditPhoto
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Field is Empty")
		return
	}

	_, payload, err := authorized(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if payload.UserId == req.UserId && req.Id != "" && req.Caption != "" {

		photo := &models.ExpUserProfile{}

		dbcontext.Table("exp_user_profile").First(photo, "id = ?", req.Id)

		photo.Caption = req.Caption
		photo.UpdatedAt = time.Now()

		err := dbcontext.Table("exp_user_profile").Save(photo).Error
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
			Message: "Unauthorized",
			Code:    "401",
		}

		c.JSON(http.StatusUnauthorized, res)
	}
}

func PhotoDeleteController(c *gin.Context) {
	var req l.DeletePhoto
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Field is Empty")
		return
	}

	_, payload, err := authorized(c)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	if payload.UserId == req.Userid {
		photo := &models.ExpUserProfile{}

		dbcontext.Table("exp_user_profile").First(photo, "id = ?", req.Id)

		photo.Deleted = true
		photo.UpdatedAt = time.Now()

		err := dbcontext.Table("exp_user_profile").Save(photo).Error
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
			Message: "Unauthorized",
			Code:    "401",
		}

		c.JSON(http.StatusUnauthorized, res)
	}
}

func Paginate(param l.GetListPhotos) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		page, _ := strconv.Atoi(param.Page)
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(param.PageSize)
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
