package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Daniel-Burbridge-Developer/Chirpy/models"
)

func main() {
	apiCfg := &models.ApiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.RequestCountHandler)
	mux.HandleFunc("/api/reset", apiCfg.ResetCountHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	corsMux := middlewareCors(mux)
	httpServer := &http.Server{Addr: "localhost:8080", Handler: corsMux}

	httpServer.ListenAndServe()
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// Not handling cases where there is no Json.Body at all. This just returns "valid", I'm fairly sure I've written this in a very hacky, not proper way.
func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnValsValid struct {
		Valid bool `json:"valid"`
	}

	type returnValsInvalid struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	fmt.Printf("Request literal: %v\n", params)
	fmt.Printf("Character length of request: %v\n", len(params.Body))

	// This seems hacky and not how I am meant to be doing this, pretty sure I'm meant to be doing this with the error value....
	if err != nil || len(params.Body) > 140 {

		respBody := returnValsInvalid{
			Error: "Something went wrong",
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	respBody := returnValsValid{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
