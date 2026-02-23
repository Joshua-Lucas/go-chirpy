package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

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
	}

	mux := http.NewServeMux()

	apiCfg.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
