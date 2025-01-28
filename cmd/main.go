package main

import (
	"forum/internal/models"
	"forum/internal/ratelimiter"
	"forum/internal/sql"
	"log"
	"net/http"
)

func main() {
	db, err := sql.InitDB()
	if err != nil {
		log.Fatal(err, "server shutdown")
	}
	models.SetDB(db)

	mux := http.NewServeMux()
	SetupRoutes(mux)

	limitedMux := ratelimiter.RateLimitMiddleware(mux)

	// Start the server
	StartServer(limitedMux)
}
