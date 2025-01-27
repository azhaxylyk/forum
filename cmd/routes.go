package main

import (
	"forum/internal/handlers"
	"net/http"
)

func SetupRoutes() {
	setupAdminRoutes()
	setupAuthRoutes()
	setupPostRoutes()
	setupNotificationRoutes()
	setupStaticRoutes()
}

func setupAdminRoutes() {
	http.HandleFunc("/admin", handlers.AdminPageHandler())
	http.HandleFunc("/admin/handle_request", handlers.HandleModerationRequest())
	http.HandleFunc("/admin/approve_moderator", handlers.AdminApproveModeratorHandler)
	http.HandleFunc("/admin/demote_moderator", handlers.DemoteModeratorHandler)
}

func setupAuthRoutes() {
	http.HandleFunc("/auth/google", handlers.GoogleAuthHandler)
	http.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)
	http.HandleFunc("/auth/github", handlers.GitHubAuthHandler)
	http.HandleFunc("/auth/github/callback", handlers.GitHubCallbackHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
}

func setupPostRoutes() {
	http.HandleFunc("/", handlers.MainPageHandler)
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
	http.HandleFunc("/profile", handlers.ProfilePageHandler)
}

func setupNotificationRoutes() {
	http.HandleFunc("/notifications", handlers.GetNotificationsHandler)
	http.HandleFunc("/notifications/mark-as-read", handlers.MarkNotificationAsReadHandler)
	http.HandleFunc("/notifications/mark-all-as-read", handlers.MarkAllNotificationsAsReadHandler)
	http.HandleFunc("/notifications/unread-count", handlers.GetUnreadCountHandler)
}

func setupStaticRoutes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("./web/icons"))))
}
