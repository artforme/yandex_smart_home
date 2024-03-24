package authTokenGenerator

import (
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"yandex_smart_house/internal/tokenApi"
)

type Response struct {
	Status uint `json:"status"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Error("failed to parse request form", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{Status: http.StatusBadRequest})
			return
		}

		code := r.PostFormValue("code")

		_, err := tokenApi.ValidateJWTToken(code)
		if err != nil {
			log.Error("token is invalid", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{Status: http.StatusBadRequest})
			return
		}
		fmt.Println(code)

	}
}
