package main

import (
	"forum/internal/config"
	"forum/internal/handlers"
	"forum/internal/models"
	"forum/internal/ratelimiter"
	"forum/internal/sql"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sql.InitDB()
	if err != nil {
		log.Fatal(err, "server shutdown")
	}
	models.SetDB(db)
	config.LoadEnv()
	handlers.InitOAuthConfigs()
	log.Printf("Google Client ID: %s", os.Getenv("GOOGLE_CLIENT_ID"))
	log.Printf("GitHub Client ID: %s", os.Getenv("GITHUB_CLIENT_ID"))
	mux := http.NewServeMux()
	SetupRoutes(mux)

	limitedMux := ratelimiter.RateLimitMiddleware(mux)

	StartServer(limitedMux)
}
