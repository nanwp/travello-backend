package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/helper"
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
			helper.ResponseOutput(c, http.StatusUnauthorized, "UNAUTHORIZE", "Unauthorize", nil)
			c.Abort()
			return
		case jwt.ValidationErrorExpired:
			helper.ResponseOutput(c, http.StatusUnauthorized, "UNAUTHORIZE", "Unauthorize! Token Expired", nil)
			c.Abort()
			return
		default:
			helper.ResponseOutput(c, http.StatusUnauthorized, "UNAUTHORIZE", "Unauthorize", nil)
			c.Abort()
			return
		}

	}

	if !token.Valid {
		helper.ResponseOutput(c, http.StatusUnauthorized, "UNAUTHORIZE", "Unauthorize", nil)
		c.Abort()
		return
	}

	UserID = claims.UserID
	c.Next()
}
