package handlers

import (
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
)

func ProfilePageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		ErrorHandler(w, r, http.StatusUnauthorized, "Please log in to view your profile")
		return
	}

	sessionToken := cookie.Value
	userID, username, err := models.GetIDBySessionToken(sessionToken)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error retrieving user session")
		return
	}

	userRole, err := models.GetUserRole(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error retrieving user role")
		return
	}

	section := r.URL.Query().Get("section")

	var (
		posts         []models.Post
		likedPosts    []models.Post
		dislikedPosts []models.Post
		comments      []map[string]interface{}
	)

	switch section {
	case "comments":
		comments, err = models.GetCommentsByUser(userID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching user comments")
			return
		}
	case "likes":
		likedPosts, err = models.GetLikedPostsByUser(userID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching liked posts")
			return
		}
	case "dislikes":
		dislikedPosts, err = models.GetDislikedPostsByUser(userID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching disliked posts")
			return
		}
	default:
		posts, err = models.GetPostsByUser(userID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching user posts")
			return
		}
	}

	tmpl, err := template.ParseFiles("web/templates/profile.html")
	if err != nil {
		log.Printf("Template loading error: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError, "Error loading template")
		return
	}

	data := struct {
		LoggedIn      bool
		Username      string
		Section       string
		Posts         []models.Post
		LikedPosts    []models.Post
		DislikedPosts []models.Post
		Comments      []map[string]interface{}
		Role          string
	}{
		LoggedIn:      true,
		Username:      username,
		Section:       section,
		Posts:         posts,
		LikedPosts:    likedPosts,
		DislikedPosts: dislikedPosts,
		Comments:      comments,
		Role:          userRole,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
	}
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
