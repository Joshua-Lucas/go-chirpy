package main

import (
	"net/http"
)

func main() {

	server := http.Server{}

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	server.Handler = mux

	server.Addr = ":8080"

	server.ListenAndServe()

}
