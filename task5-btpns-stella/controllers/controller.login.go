package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	l "task5-btpns-stella/app"
	db "task5-btpns-stella/database"
	helper "task5-btpns-stella/helpers"
	middle "task5-btpns-stella/middlewares"
	models "task5-btpns-stella/models"
	"time"
)

var (
	newReq = func(method, url string, body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, url, body)
	}
	clientDo = func(c *http.Client, req *http.Request) (*http.Response, error) {
		return c.Do(req)
	}
	dbcontext = db.ConnectDB()
)

// @BasePath /api/v1

// Login godoc
// @Summary Login Service
// @Schemes
// @Description Do Login
// @Tags Login
// @Accept json
// @Produce json
// @Param username query string false  "string valid"       minlength(5)  maxlength(10)
// @Param password query string false  "string valid"       minlength(3)  maxlength(10)
// @Success 200 {object} models.SuccessResponse
// @Router /user/login [get]
func LoginController(c *gin.Context) {
	var req l.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "Sediakan username dan password")
		return
	}

	err := helper.ValidateUser(req)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	checkUser := dbcontext.Table("exp_user eu").
		Select("COUNT(*)").
		Where("eu.username = ? AND eu.active = true", req.Username)

	var total int64
	checkUser.Count(&total)

	if total > 0 {
		getUser := dbcontext.Table("exp_user eu").
			Select("eu.id, eu.active").
			Where("eu.username = ? AND eu.active = true", req.Username)

		var userId string
		var isActive bool
		err = getUser.Limit(1).Row().Scan(&userId, &isActive)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Querying User : "+err.Error())
			return
		}

		if !isActive {
			res := l.ApiResponse{
				Code:    "200",
				Message: "User ditemukan berstatus tidak aktif",
			}

			c.JSON(http.StatusOK, res)
		}

		keycloakURL := os.Getenv("KEYCLOAK_URL") + "/realms/" + os.Getenv("KEYCLOAK_REALM_NAME") + "/protocol/openid-connect/token"
		clientToken := base64.StdEncoding.EncodeToString([]byte(os.Getenv("KEYCLOAK_CLIENT_SECRET")))
		form := url.Values{
			"grant_type": {"password"},
			"username":   {req.Username},
			"password":   {req.Password},
			"scope":      {"openid"},
		}

		client := &http.Client{Timeout: 10 * time.Second}
		request, err := newReq("POST", keycloakURL, bytes.NewBufferString(form.Encode()))
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Keycloak Error : "+err.Error())
			return
		}
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Authorization", "Basic "+clientToken)
		resp, err := clientDo(client, request)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Keycloak Request Error : "+err.Error())
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			if resp.StatusCode == http.StatusUnauthorized {
				errorResponse(c, http.StatusUnauthorized, "Status Unauthorized")
				return
			}
			errorResponse(c, http.StatusInternalServerError, "Keycloak Internal Error")
			return
		}
		respByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Error Response : "+err.Error())
			return
		}
		signIn := l.SignInResponse{}
		if err := json.Unmarshal(respByte, &signIn); err != nil {
			errorResponse(c, http.StatusInternalServerError, "Unmarshalling : "+err.Error())
			return
		}

		err = insertSession(userId)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Insert Session : "+err.Error())
			return
		}

		key, err := saveToken(userId, req.Username, signIn.AccessToken)
		if err != nil {
			errorResponse(c, http.StatusInternalServerError, "Insert to Redis : "+err.Error())
			return
		}

		res := models.SuccessResponse{
			Username: key,
			Token:    signIn.AccessToken,
		}

		c.JSON(http.StatusOK, res)

	} else {
		errorResponse(c, http.StatusInternalServerError, "User Not Found")
		return
	}

}

func saveToken(userID, username, token string) (string, error) {
	t, _ := strconv.Atoi("120")
	exp := time.Duration(t) * time.Minute

	redisVal := &middle.RedisValue{
		Username: username,
		UserId:   userID,
	}

	body, err := json.Marshal(redisVal)
	if err != nil {
		return "", err
	}

	_, err = middle.WriteRedis(token, string(body), exp)
	if err != nil {
		return "", err
	}
	return username, nil
}

func insertSession(userId string) error {

	var exists bool
	err := dbcontext.Table("exp_user_session").
		Select("count(*) > 0").
		Where("exp_user_id = ?", userId).
		Find(&exists).Error
	if err != nil {
		return err
	}

	if exists {
		err := dbcontext.Table("exp_user_session").
			Where("exp_user_id = ?", userId).
			Updates(map[string]interface{}{"last_login": "now()", "updated_at": "now()"}).
			Error
		if err != nil {
			return err
		}
	} else {
		userSession := &models.ExpUserSession{
			ExpUserId: userId,
			LastLogin: time.Now(),
			CreatedAt: time.Now(),
		}

		err := dbcontext.Table("exp_user_session").Create(userSession).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"message": message})
}
