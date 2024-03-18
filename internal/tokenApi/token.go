package tokenApi

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 7 * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("CLIENT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func RedirectToProvider(redirectURI string, responseType string, token string, state string, client_id string, scope string) string {
	return fmt.Sprintf("%s?%s=%s&client_id=%s&state=%s&scope=%s", redirectURI, responseType, token, client_id, state, scope)
}
