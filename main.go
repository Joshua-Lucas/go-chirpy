package main

import (
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/Joshua-Lucas/go-chirpy/internal/database"
	"github.com/Joshua-Lucas/go-chirpy/internal/server"
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

	dbQueries := database.New(db)

	apiCfg := server.APIConfig{
		FileserverHits: atomic.Int32{},
		DBQueries:      dbQueries,
	}

	srv := server.NewServer(&apiCfg)

	log.Fatal(srv.ListenAndServe(":8080"))

}
