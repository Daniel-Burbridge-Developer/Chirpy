package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	//mux.Handle("/app/app", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/", http.FileServer(http.Dir("./app")))
	mux.HandleFunc("/healthz", readinessHandler)

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
