package handlers

import (
	"database/sql"
	"fmt"
	"forum/internal/models"
	"html/template"
	"net/http"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	content := r.FormValue("content")
	categories := r.Form["categories"]

	content = models.SanitizeInput(content)
	if !models.IsValidContent(content) || len(categories) == 0 {
		ErrorHandler(w, r, http.StatusBadRequest, "Content and at least one category are required to create a post")
		return
	}

	postID, err := models.CreatePost(userID, content)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error creating post")
		return
	}

	for _, categoryID := range categories {
		err = models.AddCategoryToPost(postID, categoryID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error associating category")
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postID := r.FormValue("post_id")

	err = models.LikePost(userID, postID)
	if err != nil {
		http.Error(w, "Error liking post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.UpdatePostLikesDislikes(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating like count")
		return
	}

	// Уведомление владельца поста
	postOwnerID, err := models.GetPostOwner(postID)
	if err == nil && postOwnerID != userID {
		_ = models.CreateNotification(postOwnerID, userID, "like", postID, "post")
	}

	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postID := r.FormValue("post_id")

	err = models.DislikePost(userID, postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error disliking post")
		return
	}

	err = models.UpdatePostLikesDislikes(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating dislike count")
		return
	}

	// Уведомление владельца поста
	postOwnerID, err := models.GetPostOwner(postID)
	if err == nil && postOwnerID != userID {
		_ = models.CreateNotification(postOwnerID, userID, "dislike", postID, "post")
	}

	referer := r.Header.Get("Referer")
	http.Redirect(w, r, referer, http.StatusSeeOther)
}

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Missing post ID")
		return
	}

	post, err := models.GetPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorHandler(w, r, http.StatusNotFound, "Post not found")
			return
		}
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching post")
		return
	}

	comments, err := models.GetCommentsForPost(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching comments")
		return
	}

	notification := r.URL.Query().Get("notification")

	tmpl, err := template.ParseFiles("web/templates/comments.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	var loggedIn bool
	var isModerator bool
	var username string

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie == nil || cookie.Value == "" {
		loggedIn = false
	} else {
		userID, username, err := models.GetIDBySessionToken(cookie.Value)
		if err != nil || userID == "" || username == "" {
			loggedIn = false
		} else {
			loggedIn = true
			isModerator = IsModerator(userID)
		}
	}

	data := struct {
		Post         models.Post
		Comments     []models.Comment
		LoggedIn     bool
		Username     string
		IsModerator  bool // Добавьте это поле
		Notification string
	}{
		Post:         post,
		Comments:     comments,
		LoggedIn:     loggedIn,
		Username:     username,
		IsModerator:  isModerator, // По умолчанию false, если пользователь не модератор
		Notification: notification,
	}
	fmt.Println(isModerator)
	tmpl.Execute(w, data)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	postID := r.FormValue("post_id")
	ownerID, err := models.GetPostOwner(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching post owner")
		return
	}

	// Проверяем, является ли пользователь владельцем поста
	if ownerID != userID {
		ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to delete this post")
		return
	}

	err = models.DeletePost(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error deleting post")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		fmt.Println('d')
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	postID := r.FormValue("post_id")
	newContent := r.FormValue("content")
	fmt.Println(postID, newContent)

	ownerID, err := models.GetPostOwner(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching post owner")
		return
	}

	// Проверяем, является ли пользователь владельцем поста
	if ownerID != userID {
		ErrorHandler(w, r, http.StatusForbidden, "You are not allowed to edit this post")
		return
	}

	err = models.UpdatePost(postID, newContent)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating post")
		return
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}
