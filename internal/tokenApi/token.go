package tokenApi

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Second * 3).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func ValidateJWTToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("CLIENT_SECRET")), nil
	})
}
