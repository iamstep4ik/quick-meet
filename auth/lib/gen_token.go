package lib

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/iamstep4ik/quick-meet/models"
)

const (
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	Id int `json:"id"`
}

func GenerateToken(user *models.User) (string, error) {
	claims := tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Id: user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := os.Getenv("SIGNING_KEY")

	return token.SignedString([]byte(key))
}
