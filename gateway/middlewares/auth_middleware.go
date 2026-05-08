package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/gateway/grpc_client"
	"github.com/ThuraMinThein/gateway/pkg/cache"
	"github.com/ThuraMinThein/gateway/pkg/helper"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		tokenType := strings.Split(token, " ")

		if token == "" || len(tokenType) != 2 || tokenType[0] != "Bearer" {
			abortError(c, http.StatusUnauthorized)
			return
		}

		jwt := tokenType[1]

		claims, err := helper.ParseToken(jwt)

		if err != nil {
			abortError(c, http.StatusUnauthorized)
			return
		}

		var user *api.User

		if err = cache.GetCacheData("auth:"+claims.Sub, &user); err != nil {
			userClient := grpc_client.GetUserClient()
			user, err := userClient.GetUserById(c, &api.GetUserByIdRequest{UserId: claims.Sub})
			if err != nil {
				abortError(c, http.StatusUnauthorized)
				return
			}
			cache.SetCacheData("auth:"+claims.Sub, user, time.Duration(60)*time.Minute)
		}

		if err != nil {
			abortError(c, http.StatusUnauthorized)
			return
		}
		c.Set("user_id", user.Id)
		c.Next()

	}
}

func abortError(c *gin.Context, status int, message ...string) {
	errorMessage := ""
	switch status {
	case http.StatusUnauthorized:
		errorMessage = "Unauthorized"
	case http.StatusForbidden:
		errorMessage = "Forbidden"
	default:
		errorMessage = "Error"
	}
	if len(message) > 0 {
		errorMessage = errorMessage + ": " + message[0]
	}
	c.JSON(status, gin.H{"error": errorMessage})
	c.Abort()
}
