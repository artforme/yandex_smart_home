package login

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Response
		log.Info("i'm here")
		// decode body request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			panic("ouch")
		}
		render.JSON(w, r, Response{Email: req.Email, Password: req.Password})
	}
}
