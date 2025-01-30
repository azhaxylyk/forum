package main

import (
	"forum/internal/handlers"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	setupAdminRoutes(mux)
	setupAuthRoutes(mux)
	setupPostRoutes(mux)
	setupNotificationRoutes(mux)
	setupStaticRoutes(mux)
}

func setupAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin", handlers.AdminPageHandler())
	mux.HandleFunc("/admin/handle_request", handlers.HandleModerationRequest())
	mux.HandleFunc("/admin/approve_moderator", handlers.AdminApproveModeratorHandler)
	mux.HandleFunc("/admin/demote_moderator", handlers.DemoteModeratorHandler)
	mux.HandleFunc("/request_deletion", handlers.RequestDeletionHandler())
	mux.HandleFunc("/request-moderator", handlers.RequestModeratorHandler)
	mux.HandleFunc("/admin/add_category", handlers.AddCategoryHandler())
}

func setupAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/google", handlers.GoogleAuthHandler)
	mux.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)
	mux.HandleFunc("/auth/github", handlers.GitHubAuthHandler)
	mux.HandleFunc("/auth/github/callback", handlers.GitHubCallbackHandler)
	mux.HandleFunc("/register", handlers.RegisterHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)
}

func setupPostRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.MainPageHandler)
	mux.HandleFunc("/create_post", handlers.CreatePostHandler)
	mux.HandleFunc("/post", handlers.PostPageHandler)
	mux.HandleFunc("/delete_post", handlers.DeletePostHandler)
	mux.HandleFunc("/edit_post", handlers.EditPostHandler)
	mux.HandleFunc("/like", handlers.LikeHandler)
	mux.HandleFunc("/dislike", handlers.DislikeHandler)
	mux.HandleFunc("/create_comment", handlers.CreateCommentHandler)
	mux.HandleFunc("/like_comment", handlers.LikeCommentHandler)
	mux.HandleFunc("/dislike_comment", handlers.DislikeCommentHandler)
	mux.HandleFunc("/my_posts", handlers.MyPostsHandler)
	mux.HandleFunc("/liked_posts", handlers.LikedPostsHandler)
	mux.HandleFunc("/profile", handlers.ProfilePageHandler)
}

func setupNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/notifications", handlers.GetNotificationsHandler)
	mux.HandleFunc("/notifications/mark-as-read", handlers.MarkNotificationAsReadHandler)
	mux.HandleFunc("/notifications/mark-all-as-read", handlers.MarkAllNotificationsAsReadHandler)
	mux.HandleFunc("/notifications/unread-count", handlers.GetUnreadCountHandler)
}

func setupStaticRoutes(mux *http.ServeMux) {
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	mux.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("./web/icons"))))
}
