package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/natha/kyoyu-backend/initializers"
	"github.com/natha/kyoyu-backend/models"
)

func RequireAuth(c *gin.Context, role models.UserRole) {
	// Get the cookie off req
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
	}

	// Decode/validate
	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token has expired",
			})
		}

		// Find user with token
		subClaim, ok := claims["sub"].(string)
		if !ok {
			// Handle the case where "sub" claim is not a valid string
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Sub not valid",
			})
			return
		}

		// Find user by ID
		var user models.User
		initializers.DB.Where("user_id = ?", subClaim).First(&user)

		// Check if there is required roles
		if role != "" {
			if user.Role != role {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized",
				})
				return
			}
		}

		// Attach to req
		c.Set("user", user)

		// Continue
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
	}

}
