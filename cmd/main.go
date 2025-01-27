package main

import (
	"forum/internal/models"
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
	SetupRoutes()

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	log.Println("Server started on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
