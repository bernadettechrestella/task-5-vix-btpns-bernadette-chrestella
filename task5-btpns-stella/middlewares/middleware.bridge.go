package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "You need to be authorized to access this api"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func TokenValid(r *http.Request) error {
	authorizationHeader := r.Header.Get("Authorization")
	if !strings.Contains(authorizationHeader, "Bearer") {
		return errors.New("invalid token")
	}
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
	payloadMessage, err := ReadRedis(tokenString)
	if err != nil {
		return err
	}
	if payloadMessage == "" {
		return errors.New("invalid redis payload")
	}
	return nil
}
