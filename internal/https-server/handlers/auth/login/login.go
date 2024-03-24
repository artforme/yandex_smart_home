package login

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"yandex_smart_house/internal/tokenApi"
)

type Response struct {
	Status      uint   `json:"status"`
	RedirectURL string `json:"redirectURL"`
}
type Request struct {
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ResponseType string `json:"response_type"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
	UserID       string `json:"userid"`
	Password     string `json:"password"`
}

type Checker interface {
	Search(string, string) error
}

func New(log *slog.Logger, checker Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req Request
		// decode body request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode JSON request body", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{Status: http.StatusBadRequest})
			return
		}
		if err = checker.Search(req.UserID, req.Password); err != nil {
			log.Error("failed to find outgoing userID", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			if errors.Is(err, sql.ErrNoRows) {
				render.JSON(w, r, Response{Status: http.StatusNotFound})
			}
			render.JSON(w, r, Response{Status: http.StatusInternalServerError})
			return
		}

		token, err := tokenApi.GenerateToken()
		if err != nil {
			log.Error("failed to generate token", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, Response{Status: http.StatusInternalServerError})
			return
		}

		redirectURL := fmt.Sprintf("%s?%s=%s&client_id=%s&state=%s&scope=%s", req.RedirectURI, req.ResponseType, token, req.ClientID, req.State, req.Scope)

		render.JSON(w, r, Response{Status: http.StatusFound, RedirectURL: redirectURL})
	}
}
