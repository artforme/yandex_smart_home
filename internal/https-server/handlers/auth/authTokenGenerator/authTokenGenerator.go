package authTokenGenerator

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Метод запроса:", r.Method)
		fmt.Println("URI запроса:", r.RequestURI)
		fmt.Println("Хост запроса:", r.Host)
		fmt.Println("Заголовки запроса:")
		for name, headers := range r.Header {
			for _, h := range headers {
				fmt.Printf("\t%s: %s\n", name, h)
			}
		}

		// Чтение тела запроса
		data, _ := io.ReadAll(r.Body)
		fmt.Println("Тело запроса:", string(data))

		// Вывод параметров запроса (query parameters)
		fmt.Println("Параметры запроса:")
		r.ParseForm()
		for key, values := range r.Form {
			for _, value := range values {
				fmt.Printf("\t%s: %s\n", key, value)
			}
		}

	}
}
