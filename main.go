package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Joshua-Lucas/go-chirpy/internal/api"
	"github.com/Joshua-Lucas/go-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	// DB Connection
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := api.APIConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      database.New(db),
		Platform:       os.Getenv("PLATFORM"),
		TokenSecert:    os.Getenv("JWT_SIGNING_KEY"),
	}

	mux := http.NewServeMux()

	apiCfg.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
