package main

import (
	"context"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yandex_smart_house/internal/config"
	"yandex_smart_house/internal/https-server/handlers/auth/authrizetor"
	"yandex_smart_house/internal/https-server/handlers/auth/login"
	"yandex_smart_house/internal/https-server/handlers/checkAccessibility"
	"yandex_smart_house/internal/https-server/handlers/checkChangingDevices"
	"yandex_smart_house/internal/https-server/handlers/checkDeviceStatus"
	"yandex_smart_house/internal/https-server/handlers/checkListUpdate"
	"yandex_smart_house/internal/https-server/handlers/checkUserDisconnection"
	"yandex_smart_house/internal/https-server/server"
	"yandex_smart_house/internal/logger"
	"yandex_smart_house/internal/storage/postgres"
)

func main() {

	conf := config.MustLoad()

	log := logger.SetupLogger(conf.Env)

	log.Debug("debug message enabled")

	log.Info("info message enabled")

	log.Warn("warn message enabled")

	log.Error("error message enabled")

	log.Info("", slog.StringValue(conf.Address))

	storage, err := postgres.New()

	if err != nil {
		log.Error("failed to create storage", slog.StringValue(err.Error()))
		os.Exit(1)
	}

	router := mux.NewRouter()

	// There are some basic handlers for yandex that we use in specific method and url
	router.HandleFunc("/v1.0", checkAccessibility.New(log)).Methods("GET")
	router.HandleFunc("/v1.0/user/unlink", checkUserDisconnection.New(log)).Methods("POST")
	router.HandleFunc("/v1.0/user/devices", checkListUpdate.New(log)).Methods("GET")
	router.HandleFunc("/v1.0/user/devices/query", checkDeviceStatus.New(log)).Methods("POST")
	router.HandleFunc("/v1.0/user/devices/action", checkChangingDevices.New(log)).Methods("POST")
	router.HandleFunc("/api/auth/authorize", authrizetor.New(log)).Methods("GET")
	router.HandleFunc("/api/signup", login.New(log, storage)).Methods("POST")
	router.HandleFunc("/api/login", login.New(log, storage)).Methods("POST")
	//router.HandleFunc("/api/login", authrizetor.New(log)).Methods("POST")

	// setup server
	srv := &http.Server{
		Addr:         conf.Address,
		Handler:      router,
		ReadTimeout:  conf.HTTPServer.Timeout,
		WriteTimeout: conf.HTTPServer.Timeout,
		IdleTimeout:  conf.HTTPServer.Idle_timeout,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.StringValue(err.Error()))
			os.Exit(1)
		}
	}()
	log.Info("Server started")

	<-stopChan
	log.Info("Shutdown signal received, exiting...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown failed", slog.StringValue(err.Error()))
	}

	if err = server.CloseDataBase(log, storage); err != nil {
		log.Error("failed to close resources on shutdown", slog.StringValue(err.Error()))
	}

	log.Info("Server gracefully stopped")

	// TODO: make simple oauth 2.0 verification

	// TODO: test system

	// TODO: init rest api service for yandex smart house

	// TODO: deploy project

	// TODO: test project
}

//6IyGSHZlOZMLCjGPrtervr8S5wV15GeN
