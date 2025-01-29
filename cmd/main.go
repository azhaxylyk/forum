package main

import (
	"forum/internal/config"
	"forum/internal/handlers"
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
	config.LoadEnv()
	handlers.InitOAuthConfigs()
	mux := http.NewServeMux()
	SetupRoutes(mux)

	limitedMux := ratelimiter.RateLimitMiddleware(mux)

	StartServer(limitedMux)
}
