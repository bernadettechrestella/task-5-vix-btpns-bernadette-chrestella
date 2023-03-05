package controllers

import (
	"context"
	gocloak "github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	l "task5-btpns-stella/app"
	"task5-btpns-stella/models"
	"time"
)

// @BasePath /api/v1

// Register godoc
// @Summary Register New User Service
// @Schemes
// @Description Do Registration
// @Tags Register
// @Accept mpfd
// @Produce json
// @Success 200 {string} Success
// @Router /user/register [post]
func RegisterController(c *gin.Context) {
	var req l.Register

	req.Email = c.PostForm("email")
	req.Password = c.PostForm("password")
	req.Username = c.PostForm("username")

	f, err := c.FormFile("photo")
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Status Unauthorized")
		return
	}

	//isAuthorized, err := authorized(c)
	//if err != nil {
	//	errorResponse(c, http.StatusUnauthorized, err.Error())
	//	return
	//}

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

	user := &models.ExpUser{
		Username:  req.Username,
		Email:     req.Email,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = createUserKeycloak(user, req.Password)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, errs.Error())
		return
	}

	err = dbcontext.Table("exp_user").
		Create(user).Error
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, errs.Error())
		return
	}

	photo := &models.ExpUserProfile{
		ExpUserId: user.Id,
		PhotoUrl:  bucket + "/" + f.Filename,
		Title:     "profile-photo",
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

}

func createUserKeycloak(param *models.ExpUser, password string) error {

	token, err := client.LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		return err
	}

	user := gocloak.User{
		FirstName: gocloak.StringP(param.Username),
		Email:     gocloak.StringP(param.Email),
		Enabled:   gocloak.BoolP(true),
		Username:  gocloak.StringP(param.Username),
	}

	userIdKeyCloak, err := client.CreateUser(ctx, token.AccessToken, "master", user)
	if err != nil {
		return err
	}

	err = client.SetPassword(ctx, token.AccessToken, userIdKeyCloak, "master", password, false)
	if err != nil {
		return err
	}

	return nil
}
