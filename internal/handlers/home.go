package handlers

import (
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	var username string
	loggedIn := false
	var userID string

	cookie, err := r.Cookie("session_token")
	if err == nil {
		sessionToken := cookie.Value
		userID, username, err = models.GetIDBySessionToken(sessionToken)
		if err == nil {
			loggedIn = true
		}
	}

	categoryID := r.URL.Query().Get("category")
	posts, err := models.GetFilteredPosts(loggedIn, userID, categoryID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching posts")
		return
	}

	categories, err := models.GetAllCategories()
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching categories")
		return
	}

	var notifications []models.Notification
	var formattedNotifications []struct {
		models.Notification
		Message string
	}
	unreadCount := 0
	if loggedIn {
		notifications, err = models.GetNotificationsForUser(userID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error fetching notifications")
			return
		}

		for _, n := range notifications {
			message, err := n.GetMessage()
			if err != nil {
				log.Println("Error generating notification message:", err)
				message = "Error generating message"
			}
			formattedNotifications = append(formattedNotifications, struct {
				models.Notification
				Message string
			}{
				Notification: n,
				Message:      message,
			})
		}

		unreadCount, err = models.GetUnreadNotificationCount(userID)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error retrieving unread notification count")
			return
		}
	}

	notification := r.URL.Query().Get("notification")

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error loading template")
		return
	}
	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	data := struct {
		Posts         []models.Post
		Categories    []models.Category
		LoggedIn      bool
		Username      string
		Notifications []struct {
			models.Notification
			Message string
		}
		Notification             string
		SelectedCategory         string
		UnreadNotificationsCount int
	}{
		Posts:                    posts,
		Categories:               categories,
		LoggedIn:                 loggedIn,
		Username:                 username,
		Notifications:            formattedNotifications,
		Notification:             notification,
		SelectedCategory:         categoryID,
		UnreadNotificationsCount: unreadCount,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
