package handlers

import (
	"forum/internal/models"
	"net/http"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	content = models.SanitizeInput(content)
	if !models.IsValidContent(content) {
		ErrorHandler(w, r, http.StatusBadRequest, "Content is required to create a comment")
		return
	}

	err = models.CreateComment(postID, userID, content)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error creating comment")
		return
	}

	postOwnerID, err := models.GetPostOwner(postID)
	if err == nil && postOwnerID != userID {
		_ = models.CreateNotification(postOwnerID, userID, "comment", postID, "post")
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	commentID := r.FormValue("comment_id")
	postID := r.FormValue("post_id")

	err = models.LikeComment(userID, commentID)
	if err != nil {
		if err.Error() == "you have already liked this comment" {
			http.Redirect(w, r, "/post?id="+postID+"&notification=already_liked", http.StatusSeeOther)
			return
		}
		http.Error(w, "Error liking comment", http.StatusInternalServerError)
		return
	}

	err = models.UpdateCommentLikesDislikes(commentID)
	if err != nil {
		http.Error(w, "Error updating like count: "+err.Error(), http.StatusInternalServerError)
		return
	}

	commentOwnerID, err := models.GetCommentOwner(commentID)
	if err == nil && commentOwnerID != userID {
		_ = models.CreateNotification(commentOwnerID, userID, "like", commentID, "comment")
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	commentID := r.FormValue("comment_id")
	postID := r.FormValue("post_id")

	err = models.DislikeComment(userID, commentID)
	if err != nil {
		if err.Error() == "you have already disliked this comment" {
			http.Redirect(w, r, "/post?id="+postID+"&notification=already_disliked", http.StatusSeeOther)
			return
		}
		http.Error(w, "Error disliking comment", http.StatusInternalServerError)
		return
	}

	err = models.UpdateCommentLikesDislikes(commentID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	commentOwnerID, err := models.GetCommentOwner(commentID)
	if err == nil && commentOwnerID != userID {
		_ = models.CreateNotification(commentOwnerID, userID, "dislike", commentID, "comment")
	}

	http.Redirect(w, r, "/post?id="+postID, http.StatusSeeOther)
}
