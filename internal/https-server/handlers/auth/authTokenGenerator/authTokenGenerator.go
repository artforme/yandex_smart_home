package authTokenGenerator

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"yandex_smart_house/internal/tokenApi"
)

type Response struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    uint   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func ReturnAccessToken(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Error("failed to parse request form", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{})
			return
		}
		code := r.PostFormValue("code")

		token, err := tokenApi.ValidateJWTToken(code, false)
		if err != nil {
			log.Error("token is invalid", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		AccessToken, err := tokenApi.GenerateAccessToken(claims["userID"].(string))
		if err != nil {
			log.Error("failed to generate token", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{})
			return
		}
		render.JSON(w, r, Response{AccessToken: AccessToken, TokenType: "JWT", ExpiresIn: 60 * 60 * 24})
	}
}

func ReturnRefreshToken(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Error("failed to parse request form", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{})
			return
		}

		code := r.PostFormValue("code")

		token, err := tokenApi.ValidateJWTToken(code, false)
		if err != nil {
			log.Error("token is invalid", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		RefreshToken, err := tokenApi.GenerateRefreshToken(claims["userID"].(string))
		if err != nil {
			log.Error("failed to generate token", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{})
			return
		}
		render.JSON(w, r, Response{RefreshToken: RefreshToken, TokenType: "JWT", ExpiresIn: 60 * 60 * 24 * 7})

	}
}
