package server

import (
	"net/http"

	"go.uber.org/zap"
)

type server struct {
	logger *zap.Logger
}

func New(l *zap.Logger) *server {
	return &server{
		logger: l,
	}
}

func (s *server) HandlePrices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
	return
}
