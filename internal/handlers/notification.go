package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/models"
	"html/template"
	"net/http"
)

type TemplateData struct {
	Notifications []struct {
		Message             string
		ReceivedAtFormatted string
	}
	UnreadNotificationsCount int
}

func GetNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	notifications, err := models.GetNotificationsForUser(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error retrieving notifications")
		return
	}

	unreadCount, err := models.GetUnreadNotificationCount(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error retrieving unread notification count")
		return
	}

	var templateNotifications []struct {
		Message             string
		ReceivedAtFormatted string
	}
	for _, n := range notifications {
		message, err := n.GetMessage()
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, "Error generating notification message")
			return
		}
		templateNotifications = append(templateNotifications, struct {
			Message             string
			ReceivedAtFormatted string
		}{
			Message:             message,
			ReceivedAtFormatted: n.CreatedAtFormatted,
		})
	}

	data := TemplateData{
		Notifications:            templateNotifications,
		UnreadNotificationsCount: unreadCount,
	}

	tmpl, err := template.ParseFiles("web/templates/notifications.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if err = tmpl.Execute(w, data); err != nil {
		fmt.Println(err)
	}
}

func MarkNotificationAsReadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	notificationID := r.FormValue("notification_id")
	if notificationID == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Notification ID is required")
		return
	}

	err = models.MarkNotificationAsRead(notificationID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error marking notification as read")
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		ErrorHandler(w, r, http.StatusUnauthorized, "Invalid session token")
		return
	}

	err = models.DeleteReadNotifications(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error deleting read notifications")
		return
	}

	http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}

func MarkAllNotificationsAsReadHandler(w http.ResponseWriter, r *http.Request) {
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

	err = models.MarkAllNotificationsAsRead(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error marking all notifications as read")
		return
	}

	err = models.DeleteReadNotifications(userID)
	if err != nil {
		ErrorHandler(w, r, http.StatusInternalServerError, "Error deleting read notifications")
		return
	}

	http.Redirect(w, r, "/notifications", http.StatusSeeOther)
}

func GetUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	unreadCount, err := models.GetUnreadNotificationCount(userID)
	if err != nil {
		http.Error(w, "Error retrieving unread count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"unreadCount": unreadCount})
}
