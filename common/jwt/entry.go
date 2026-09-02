package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userId string, expire int64, secret string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Second * time.Duration(expire)).Unix(),
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
