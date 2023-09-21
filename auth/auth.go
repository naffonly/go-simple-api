package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	_ "github.com/golang-jwt/jwt"
	"net/http"
	"simple-api/models"
	"time"
)

const (
	USERNAME = "admin"
	PASSWORD = "admin"
	SECRET   = "secret"
)

func LoginHalder(c *gin.Context) {
	var user models.Credential
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Bad Request"})
	}
	// fmt.Println(user.Username)
	if user.Username != USERNAME {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
	} else if user.Password != PASSWORD {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Unauthorized"})
	} else {
		//token
		claim := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
			Issuer:    "test",
			IssuedAt:  time.Now().Unix(),
		}

		signin := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		token, err := signin.SignedString([]byte(SECRET))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Internal Server Error",
				"err": err.Error(),
			})
			c.Abort()
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":   "Login Success",
			"token": token,
		})
	}
}
