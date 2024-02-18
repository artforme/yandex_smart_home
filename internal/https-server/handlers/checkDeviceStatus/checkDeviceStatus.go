package checkDeviceStatus

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	Resp string `json:"resp"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.checkAccessibility.New"

		// add to log info
		log = log.With(
			slog.String("op", op),
		)
		render.JSON(w, r, Response{
			Resp: "200",
		})
	}
}
