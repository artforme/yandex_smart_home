package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"yandex_smart_house/internal/config"
	"yandex_smart_house/internal/logger"
)

func main() {

	conf := config.MustLoad()

	log := logger.SetupLogger(conf.Env)

	log.Debug("debug message enabled")

	log.Info("info message enabled")

	log.Warn("warn message enabled")

	log.Error("error message enabled")

	log.Info("fuck you", slog.StringValue(conf.Address))

	router := chi.NewRouter()

	// There are some basic handlers for yandex that we use in specific method and url
	router.Get("/v1.0", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World from /v1.0 on Go server!")
	})
	////router.Post("v1.0/user/unlink", somehandler.New(log))
	////router.Head("v1.0/user/devices", somehandler.New(log))
	////router.Post("/v1.0/user/devices/query", somehandler.New(log))
	////router.Post("/v1.0/user/devices/action", somehandler.New(log))

	// setup server
	srv := &http.Server{
		Addr:         conf.Address,
		Handler:      router,
		ReadTimeout:  conf.HTTPServer.Timeout,
		WriteTimeout: conf.HTTPServer.Timeout,
		IdleTimeout:  conf.HTTPServer.Idle_timeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

	// TODO: make simple oauth 2.0 verification

	// TODO: init storage and connect it

	// TODO: test system

	// TODO: init rest api service for yandex smart house

	// TODO: deploy project

	// TODO: test project
}
