package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nanwp/travello/config"
)

var UserID string

func JWTMiddleware(c *gin.Context) {
	headerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(headerToken, "Bearer ")

	tokenString := splitToken[1]

	claims := &config.JWTClaim{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		},
	)

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		switch v.Errors {
		case jwt.ValidationErrorSignatureInvalid:
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorize",
				"status":  "UNAUTHORIZE",
			})
			c.Abort()
			return
		case jwt.ValidationErrorExpired:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "UNAUTHORIZE",
				"message": "Unauthorize! Token Expired",
			})
			c.Abort()
			return
		default:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "UNAUTHORIZE",
				"message": "Unauthorize",
			})
			c.Abort()
			return
		}

	}

	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "UNAUTHORIZE",
			"message": "Unauthorize",
		})
		c.Abort()
		return
	}

	UserID = claims.UserID
	c.Next()
}
