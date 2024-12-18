package handlers

import (
	"database/sql"
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
		ErrorHandler(w, r, http.StatusInternalServerError, "Error liking post")
		return
	}

	err = models.UpdatePostLikesDislikes(postID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error updating like count")
		return
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
	var username string
	cookie, _ := r.Cookie("session_token")

	_, username, err = models.GetIDBySessionToken(cookie.Value)
	if err == nil {
		loggedIn = true
	}

	data := struct {
		Post         models.Post
		Comments     []models.Comment
		LoggedIn     bool
		Username     string
		Notification string
	}{
		Post:         post,
		Comments:     comments,
		LoggedIn:     loggedIn,
		Username:     username,
		Notification: notification,
	}

	tmpl.Execute(w, data)
}

func MyPostsHandler(w http.ResponseWriter, r *http.Request) {
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

	posts, err := models.GetPostsByUser(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching posts")
		return
	}

	categories, err := models.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching categories")
		return
	}

	tmpl, err := template.ParseFiles("web/templates/posts.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts            []models.Post
		Categories       []models.Category
		LoggedIn         bool
		Username         string
		SelectedCategory string
		SelectedFilter   string
	}{
		Posts:            posts,
		Categories:       categories,
		LoggedIn:         true,
		SelectedCategory: "",
		SelectedFilter:   "",
	}

	tmpl.Execute(w, data)
}

func LikedPostsHandler(w http.ResponseWriter, r *http.Request) {
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

	posts, err := models.GetLikedPostsByUser(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching liked posts")
		return
	}

	categories, err := models.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching categories")
		return
	}

	tmpl, err := template.ParseFiles("web/templates/posts.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Posts            []models.Post
		Categories       []models.Category
		LoggedIn         bool
		Username         string
		SelectedCategory string
		SelectedFilter   string
	}{
		Posts:            posts,
		Categories:       categories,
		LoggedIn:         true,
		SelectedCategory: "",
		SelectedFilter:   "",
	}

	tmpl.Execute(w, data)
}
