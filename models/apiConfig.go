package models

import (
	"fmt"
	"net/http"
)

type ApiConfig struct {
	fileserverHits int
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) RequestCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits)))
}

func (cfg *ApiConfig) ResetCountHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}
