package main

import (
	"net/http"
)

func main() {

	server := http.Server{}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server.Handler = mux

	server.Addr = ":8080"

	server.ListenAndServe()

}
