package main

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"yandex_smart_house/internal/config"
	"yandex_smart_house/internal/handlers/checkAccessibility"
	"yandex_smart_house/internal/handlers/checkChangingDevices"
	"yandex_smart_house/internal/handlers/checkDeviceStatus"
	"yandex_smart_house/internal/handlers/checkListUpdate"
	"yandex_smart_house/internal/handlers/checkUserDisconnection"
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

	router := mux.NewRouter()

	// There are some basic handlers for yandex that we use in specific method and url
	router.HandleFunc("/v1.0", checkAccessibility.New(log)).Methods("GET")
	router.HandleFunc("/v1.0/user/unlink", checkUserDisconnection.New(log)).Methods("POST")
	router.HandleFunc("/v1.0/user/devices", checkListUpdate.New(log)).Methods("GET")
	router.HandleFunc("/v1.0/user/devices/query", checkDeviceStatus.New(log)).Methods("POST")
	router.HandleFunc("/v1.0/user/devices/action", checkChangingDevices.New(log)).Methods("POST")

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
