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

		// Получаем запросы на модерацию
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

		// Получаем всех модераторов
		moderators, err := models.GetAllModerators()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to load moderators"))
			return
		}

		// Получаем данные о постах для каждого запроса на удаление
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

		tmpl, err := template.ParseFiles("web/templates/admin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to load template"))
			return
		}
		section := r.URL.Query().Get("section")
		// Передаем данные в шаблон
		err = tmpl.Execute(w, struct {
			DeleteRequests    []models.ModerationRequest
			ModeratorRequests []models.ModerationRequest
			Moderators        []models.User // Добавляем список модераторов
			Section           string
		}{
			DeleteRequests:    deleteRequests,
			ModeratorRequests: moderatorRequests,
			Moderators:        moderators,
			Section:           section,
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
		if !isAdmin(userID) { // Проверяем, админ ли пользователь
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Получаем ID запроса из формы
		requestID, err := strconv.Atoi(r.FormValue("request_id"))
		if err != nil {
			http.Error(w, "Invalid request ID", http.StatusBadRequest)
			return
		}

		// Получаем действие (approve или reject)
		action := r.FormValue("action")
		if action != "approve" && action != "reject" {
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}

		// Получаем данные о запросе
		request, err := models.GetModerationRequestByID(requestID)
		if err != nil {
			http.Error(w, "Failed to retrieve request", http.StatusInternalServerError)
			return
		}
		if request == nil {
			http.Error(w, "Request not found", http.StatusNotFound)
			return
		}

		// Если запрос на удаление поста и действие - approve, удаляем пост
		if action == "approve" && request.Type == "delete_post" {
			// Проверяем, существует ли значение PostID и оно валидно
			if !request.PostID.Valid {
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}

			// Удаляем пост
			err := models.DeletePostByAdmin(request.PostID.String, "Deleted by admin") // Используем String() для извлечения значения
			if err != nil {
				http.Error(w, "Failed to delete post", http.StatusInternalServerError)
				return
			}

			// Отправляем уведомление пользователю
			err = models.CreateNotification(request.UserID, userID, "approve_del", request.PostID.String, "post")
			if err != nil {
				http.Error(w, "Failed to create notification", http.StatusInternalServerError)
				return
			}
		}

		if action == "reject" && request.Type == "delete_post" {
			// Отправляем уведомление пользователю о том, что запрос был отклонен
			err := models.CreateNotification(request.UserID, userID, "reject_del", request.PostID.String, "post")
			if err != nil {
				http.Error(w, "Failed to create notification", http.StatusInternalServerError)
				return
			}
		}

		// Обновляем статус запроса
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
	// Извлекаем cookie сессии
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "" // Если cookie нет, пользователь не авторизован
	}

	// Получаем токен из cookie
	sessionToken := cookie.Value

	// Получаем ID пользователя по токену через модель
	userID, _, err := models.GetIDBySessionToken(sessionToken)
	if err != nil {
		return "" // Ошибка при извлечении ID
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

		// Сохраняем запрос на модерацию
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
	action := r.FormValue("action") // "approve" или "reject"
	requestID := r.FormValue("request_id")
	requestIDInt, _ := strconv.Atoi(requestID)

	var newRole string
	var notificationType string
	var status string
	if action == "approve" {
		newRole = "moderator"
		notificationType = "approve_mod" // Уведомление о принятии в модераторы
		status = "approved"
	} else {
		newRole = "user"
		notificationType = "reject_mod" // Уведомление об отклонении
		status = "rejected"
	}

	// Обновляем роль пользователя
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

	// Отправляем уведомление о статусе модератора
	err := models.CreateNotification(userID, userID, notificationType, "", "role")
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to send notification")
		return
	}

	// Перенаправляем на страницу админ-панели
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func DemoteModeratorHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")

	// Проверяем, есть ли права администратора
	adminID := getUserIDFromSession(r)
	if !isAdmin(adminID) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Меняем роль на "user"
	if err := models.UpdateUserRole(userID, "user"); err != nil {
		log.Printf("Error demoting moderator: %v", err)
		http.Error(w, "Failed to demote moderator", http.StatusInternalServerError)
		return
	}

	// Перенаправляем обратно на страницу админа
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
