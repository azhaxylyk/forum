package main

import (
	"forum/internal/handlers"
	"forum/internal/models"
	"forum/internal/sql"
	"log"
	"net/http"
)

func main() {
	// Инициализация базы данных
	db, err := sql.InitDB()
	if err != nil {
		log.Fatal(err, "server shutdown")
	}

	// Установка базы данных для моделей
	models.SetDB(db)

	http.HandleFunc("/auth/google", handlers.GoogleAuthHandler)
	http.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)

	// Регистрация маршрутов
	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/create_post", handlers.CreatePostHandler)
	http.HandleFunc("/post", handlers.PostPageHandler)
	http.HandleFunc("/delete_post", handlers.DeletePostHandler)
	http.HandleFunc("/edit_post", handlers.EditPostHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/dislike", handlers.DislikeHandler)
	http.HandleFunc("/create_comment", handlers.CreateCommentHandler)
	http.HandleFunc("/like_comment", handlers.LikeCommentHandler)
	http.HandleFunc("/dislike_comment", handlers.DislikeCommentHandler)
	http.HandleFunc("/my_posts", handlers.MyPostsHandler)
	http.HandleFunc("/liked_posts", handlers.LikedPostsHandler)
	http.HandleFunc("/profile", handlers.ProfilePageHandler) // Добавлен маршрут профиля
	http.HandleFunc("/notifications", handlers.GetNotificationsHandler)
	http.HandleFunc("/notifications/mark-as-read", handlers.MarkNotificationAsReadHandler)
	http.HandleFunc("/notifications/mark-all-as-read", handlers.MarkAllNotificationsAsReadHandler)
	http.HandleFunc("/notifications/unread-count", handlers.GetUnreadCountHandler)

	// Обслуживание статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("./web/icons"))))

	// Запуск сервера
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
