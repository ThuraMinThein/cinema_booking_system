package middlewares

// import (
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// )

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		token := c.GetHeader("Authorization")
// 		tokenType := strings.Split(token, " ")

// 		if token == "" || len(tokenType) != 2 || tokenType[0] != "Bearer" {
// 			abortError(c, http.StatusUnauthorized)
// 			return
// 		}

// 		jwt := tokenType[1]

// 		claims, err := helper.ParseToken(jwt)

// 		if err != nil {
// 			abortError(c, http.StatusUnauthorized)
// 			return
// 		}

// 		var user *models.User

// 		if err = cache.GetCacheData("user:"+claims.Sub, &user); err != nil {
// 			authHandler := handlers.NewAuthHandler(database.DB)
// 			user, err = authHandler.GetUserById(claims.Sub)
// 			if err != nil {
// 				abortError(c, http.StatusUnauthorized)
// 				return
// 			}
// 			cache.SetCacheData("user:"+claims.Sub, user, time.Duration(60)*time.Minute)
// 		}

// 		if err != nil {
// 			abortError(c, http.StatusUnauthorized)
// 			return
// 		}
// 		c.Set("user_id", user.ID)
// 		c.Set("company_id", user.CompanyID)
// 		c.Set("role", user.Role)
// 		c.Set("user", user)
// 		c.Next()

// 	}
// }

// func abortError(c *gin.Context, status int, message ...string) {
// 	errorMessage := ""
// 	switch status {
// 	case http.StatusUnauthorized:
// 		errorMessage = "Unauthorized"
// 	case http.StatusForbidden:
// 		errorMessage = "Forbidden"
// 	default:
// 		errorMessage = "Error"
// 	}
// 	if len(message) > 0 {
// 		errorMessage = errorMessage + ": " + message[0]
// 	}
// 	c.JSON(status, gin.H{"error": errorMessage})
// 	c.Abort()
// }
