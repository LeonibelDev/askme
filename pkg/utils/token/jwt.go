package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
		"authorized": true,
		"user":       email,
	})

	secretKey := []byte(os.Getenv("SECRET"))

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}

	return tokenString, nil

}
