package handlers

import (
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func AdminPageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromSession(r)
		if !isAdmin(userID) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}

		deleteRequests, err := models.GetAllModerationRequests()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to load delete post requests"))
			return
		}

		moderatorRequests, err := models.GetModeratorRequests()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to load moderator requests"))
			return
		}

		moderators, err := models.GetAllModerators()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to load moderators"))
			return
		}

		for i, request := range deleteRequests {
			if request.PostID.Valid {
				post, err := models.GetPostByID(request.PostID.String)
				if err != nil {
					deleteRequests[i].PostContent = "Post not found"
				} else {
					deleteRequests[i].PostContent = post.Content
				}
			}
		}

		for i, request := range moderatorRequests {
			username, err := models.GetUsernameByID(request.UserID)
			if err != nil {
				moderatorRequests[i].Username = "Username not found"
			} else {
				moderatorRequests[i].Username = username
			}
		}

		var categories []models.Category
		categories, err = models.GetAllCategories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to load categories"))
			return
		}

		tmpl, err := template.ParseFiles("web/templates/admin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to load template"))
			return
		}
		section := r.URL.Query().Get("section")
		err = tmpl.Execute(w, struct {
			DeleteRequests    []models.ModerationRequest
			ModeratorRequests []models.ModerationRequest
			Moderators        []models.User
			Section           string
			Categories        []models.Category
		}{
			DeleteRequests:    deleteRequests,
			ModeratorRequests: moderatorRequests,
			Moderators:        moderators,
			Section:           section,
			Categories:        categories,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to render template"))
		}
	}
}

func HandleModerationRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromSession(r)
		if !isAdmin(userID) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		requestID, err := strconv.Atoi(r.FormValue("request_id"))
		if err != nil {
			http.Error(w, "Invalid request ID", http.StatusBadRequest)
			return
		}

		action := r.FormValue("action")
		if action != "approve" && action != "reject" {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}

		request, err := models.GetModerationRequestByID(requestID)
		if err != nil {
			http.Error(w, "Failed to retrieve request", http.StatusInternalServerError)
			return
		}
		if request == nil {
			http.Error(w, "Request not found", http.StatusNotFound)
			return
		}

		if action == "approve" && request.Type == "delete_post" {
			if !request.PostID.Valid {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}

			err := models.DeletePostByAdmin(request.PostID.String, "Deleted by admin")
			if err != nil {
				http.Error(w, "Failed to delete post", http.StatusInternalServerError)
				return
			}

			err = models.CreateNotification(request.UserID, userID, "approve_del", request.PostID.String, "post")
			if err != nil {
				http.Error(w, "Failed to create notification", http.StatusInternalServerError)
				return
			}
		}

		if action == "reject" && request.Type == "delete_post" {
			err := models.CreateNotification(request.UserID, userID, "reject_del", request.PostID.String, "post")
			if err != nil {
				http.Error(w, "Failed to create notification", http.StatusInternalServerError)
				return
			}
		}

		err = models.UpdateModerationRequestStatus(requestID, action)
		if err != nil {
			http.Error(w, "Failed to update request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request updated successfully"))
	}
}

func isAdmin(userID string) bool {
	role, err := models.GetUserRole(userID)
	if err != nil {
		return false
	}
	return role == "admin"
}

func IsModerator(userID string) bool {
	role, err := models.GetUserRole(userID)
	if err != nil {
		return false
	}
	return role == "moderator"
}

func getUserIDFromSession(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return ""
	}

	sessionToken := cookie.Value

	userID, _, err := models.GetIDBySessionToken(sessionToken)
	if err != nil {
		return ""
	}

	return userID
}

func RequestDeletionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromSession(r)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		postID := r.FormValue("post_id")
		reason := r.FormValue("reason")

		if postID == "" || reason == "" {
			http.Error(w, "Post ID and reason are required", http.StatusBadRequest)
			return
		}

		err := models.CreateModerationRequest(userID, "delete_post", reason, postID)
		if err != nil {
			http.Error(w, "Failed to create moderation request", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Moderation request submitted successfully"))
	}
}

func AdminApproveModeratorHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	action := r.FormValue("action")
	requestID := r.FormValue("request_id")
	requestIDInt, _ := strconv.Atoi(requestID)

	var newRole string
	var notificationType string
	var status string
	if action == "approve" {
		newRole = "moderator"
		notificationType = "approve_mod"
		status = "approved"
	} else {
		newRole = "user"
		notificationType = "reject_mod"
		status = "rejected"
	}

	if err := models.UpdateUserRole(userID, newRole); err != nil {
		log.Printf("Error updating user role: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to process request")
		return
	}

	if err := models.UpdateModerationRequestStatus(requestIDInt, status); err != nil {
		log.Printf("Error updating moderation request status: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to update request status")
		return
	}

	err := models.CreateNotification(userID, userID, notificationType, "", "role")
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to send notification")
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DemoteModeratorHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")

	adminID := getUserIDFromSession(r)
	if !isAdmin(adminID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := models.UpdateUserRole(userID, "user"); err != nil {
		log.Printf("Error demoting moderator: %v", err)
		http.Error(w, "Failed to demote moderator", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func RequestModeratorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized. Please log in.", http.StatusUnauthorized)
		return
	}

	userID, _, err := models.GetIDBySessionToken(cookie.Value)
	if err != nil {
		http.Error(w, "Could not retrieve user data from session.", http.StatusUnauthorized)
		return
	}
	err = models.CreateModerationRequest(userID, "moderator_request", "User requests moderator status", "")
	if err != nil {
		log.Printf("Error creating moderation request: %v", err)
		http.Error(w, "Could not create moderation request.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func AddCategoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserIDFromSession(r)
		if !isAdmin(userID) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		categoryName := r.FormValue("category_name")
		if categoryName == "" {
			http.Error(w, "Category name is required", http.StatusBadRequest)
			return
		}

		err := models.AddCategory(categoryName)
		if err != nil {
			log.Printf("Error adding category: %v", err)
			http.Error(w, "Failed to add category", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin?section=categories", http.StatusSeeOther)
	}
}
