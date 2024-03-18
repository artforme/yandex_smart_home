package login

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	Status string `json:"status"`
}
type Request struct {
	UserID   string `json:"userid"`
	Password string `json:"password"`
}

type Checker interface {
	Search(string, string) error
}

func New(log *slog.Logger, checker Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		log.Info("i'm here")
		// decode body request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			panic("ouch")
		}
		fmt.Println(req.UserID, req.Password)
		if err = checker.Search(req.UserID, req.Password); err != nil {
			log.Error("failed to find outgoing userID", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			if errors.Is(err, sql.ErrNoRows) {
				render.JSON(w, r, Response{Status: "404"})
			}
			render.JSON(w, r, Response{Status: "400"})
			return
		}

		render.JSON(w, r, Response{Status: "200"})
	}
}
