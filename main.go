package main

import "net/http"

func main() {

	server := http.Server{}

	server.Handler = http.NewServeMux()
	server.Addr = ":8080"

	server.ListenAndServe()

}
