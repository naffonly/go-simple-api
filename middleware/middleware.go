package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

const (
	SECRET = "secret"
)

func AuthValid(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"msg": "Unauthorized",
		})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, invalid := token.Method.(*jwt.SigningMethodHMAC); !invalid {
			return nil, fmt.Errorf("Invalid Token", token.Header["alg"])
		}
		return []byte(SECRET), nil
	})

	if token != nil && err == nil {
		fmt.Println("Token verified")
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "not Unauthorized",
			"err": err.Error(),
		})
		c.Abort()
	}
}
