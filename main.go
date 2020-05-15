package main

import (
	"log"
	"net/http"

	"coinbani/server"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("initializing server")

	s := server.New(logger)

	rourter := mux.NewRouter()
	rourter.HandleFunc("/cryptos/prices", s.HandlePrices).Methods(http.MethodGet)

	logger.Info("server listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", rourter))
}
