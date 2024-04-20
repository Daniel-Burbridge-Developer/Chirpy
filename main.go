package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/Daniel-Burbridge-Developer/Chirpy/models"
)

func main() {
	apiCfg := &models.ApiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.RequestCountHandler)
	mux.HandleFunc("/api/reset", apiCfg.ResetCountHandler)
	// mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("POST /api/chirps", uploadChirpHandler)
	mux.HandleFunc("GET /api/chirps/", receiveChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{id}", receiveChirpsHandler)
	// mux.HandleFunc("GET api/chirps/*", receiveByChirpIDHandler)

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

func uploadChirpHandler(w http.ResponseWriter, r *http.Request) {
	validateChirpHandler(w, r)
}

func receiveChirpsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := models.NewDB("./internal/database/database.json")
	if err != nil {
		fmt.Printf("error initializing DB: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	chirps, err := db.GetChirps()
	if err != nil {
		fmt.Printf("error retrieving chirps: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	pv := r.PathValue("id")
	id, _ := strconv.Atoi(pv)

	if pv != "" {
		// fmt.Printf("attempting to return chirp %v\n", id)
		respondWithJSON(w, 200, chirps[id-1])
		return
	}

	// fmt.Printf("return all chirps\n")
	respondWithJSON(w, 200, chirps)
}

// func receiveByChirpIDHandler(w http.ResponseWriter, r *http.Request) {
// 	r.PathValue()
// }

// Not handling cases where there is no Json.Body at all. This just returns "valid", I'm fairly sure I've written this in a very hacky, not proper way.
func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	// fmt.Printf("Request literal: %v\n", params)
	// fmt.Printf("Character length of request: %v\n", len(params.Body))

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpBody := wordReplacer(badWords, strings.Split(params.Body, " "))

	// fmt.Printf("after badword replaced, before call to createchirp%v\n", chirpBody)

	if err != nil || len(chirpBody) > 140 {
		if err != nil {
			respondWithError(w, 400, err.Error())
		} else {
			respondWithError(w, 400, "chirp too long")
		}

	} else {
		db, _ := models.NewDB("./internal/database/database.json")
		chirp, _ := db.CreateChirp(chirpBody)
		respondWithJSON(w, 201, chirp)
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVals struct {
		Error string `json:"error"`
	}

	respBody := returnVals{
		Error: msg,
	}

	dat, err := json.Marshal(respBody)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

// Expects badWords to be passed in, in lowercase.
func wordReplacer(badWords []string, usedWords []string) string {
	cleanWords := make([]string, 0, len(usedWords))
	for _, word := range usedWords {
		if slices.Contains(badWords, strings.ToLower(word)) {
			cleanWords = append(cleanWords, "****")
		} else {
			cleanWords = append(cleanWords, word)
		}
	}
	return strings.Join(cleanWords, " ")
}
