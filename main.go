package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/*", http.FileServer(http.Dir(".")))

	corsMux := middlewareCors(mux)
	readiMux := middlewareReadiness(corsMux)

	httpServer := &http.Server{Addr: "localhost:8080", Handler: readiMux}
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

func middlewareReadiness(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		next.ServeHTTP(w, r)
	})
}
