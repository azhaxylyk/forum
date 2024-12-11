package main

import (
	"forum/internal/handlers"
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

	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/create_post", handlers.CreatePostHandler)
	http.HandleFunc("/post", handlers.PostPageHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/dislike", handlers.DislikeHandler)
	http.HandleFunc("/create_comment", handlers.CreateCommentHandler)
	http.HandleFunc("/like_comment", handlers.LikeCommentHandler)
	http.HandleFunc("/dislike_comment", handlers.DislikeCommentHandler)
	http.HandleFunc("/my_posts", handlers.MyPostsHandler)
	http.HandleFunc("/liked_posts", handlers.LikedPostsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
