package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"
)

type User struct {
	ID       string
	Username string
}

type ModerationRequest struct {
	ID          int
	UserID      string
	Type        string
	Reason      string
	PostID      sql.NullString // Null, если запрос не связан с постом
	Status      string
	CreatedAt   time.Time
	PostContent template.HTML
	Username    string
}

func GetAllModerationRequests() ([]ModerationRequest, error) {
	rows, err := db.Query(`SELECT id, user_id, type, reason, post_id, status, created_at FROM moderation_requests ORDER BY created_at DESC`)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	var requests []ModerationRequest
	for rows.Next() {
		var request ModerationRequest
		err = rows.Scan(&request.ID, &request.UserID, &request.Type, &request.Reason, &request.PostID, &request.Status, &request.CreatedAt)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if len(requests) == 0 {
		fmt.Println("No requests found in the database")
	}

	return requests, nil
}

func CreateModerationRequest(userID, requestType, reason, postID string) error {
	_, err := db.Exec(
		`INSERT INTO moderation_requests (user_id, type, reason, post_id) VALUES (?, ?, ?, ?)`,
		userID, requestType, reason, postID,
	)
	if err != nil {
		log.Printf("Error inserting moderation request: %v", err)
	}
	return err
}

func UpdateModerationRequestStatus(requestID int, status string) error {
	_, err := db.Exec(`UPDATE moderation_requests SET status = $1 WHERE id = $2`, status, requestID)
	return err
}

func GetUserRole(userID string) (string, error) {
	var role string
	err := db.QueryRow("SELECT role FROM users WHERE id = ?", userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user with ID %s not found", userID)
		}
		return "", err
	}
	return role, nil
}

func DeletePostByAdmin(postID string, reason string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var userID string
	err = tx.QueryRow("SELECT user_id FROM posts WHERE id = ?", postID).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Удаление поста и связанных данных
	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM post_likes WHERE post_id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM comments WHERE post_id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM posts WHERE id = ?", postID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Фиксируем транзакцию
	return tx.Commit()
}

func GetModerationRequestByID(requestID int) (*ModerationRequest, error) {
	var request ModerationRequest

	// Выполняем запрос к базе данных для получения данных о запросе по ID
	row := db.QueryRow(`
        SELECT id, user_id, type, reason, post_id, status
        FROM moderation_requests
        WHERE id = ?`, requestID)

	err := row.Scan(&request.ID, &request.UserID, &request.Type, &request.Reason, &request.PostID, &request.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Если запрос не найден, возвращаем nil
		}
		return nil, err // В случае другой ошибки возвращаем ошибку
	}

	return &request, nil
}

func GetModeratorRequests() ([]ModerationRequest, error) {
	rows, err := db.Query("SELECT id, user_id, type, reason, post_id, status FROM moderation_requests WHERE type = 'moderator_request'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []ModerationRequest
	for rows.Next() {
		var req ModerationRequest
		err := rows.Scan(&req.ID, &req.UserID, &req.Type, &req.Reason, &req.PostID, &req.Status)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func UpdateUserRole(userID, newRole string) error {
	_, err := db.Exec("UPDATE users SET role = ? WHERE id = ?", newRole, userID)
	return err
}

func GetAllModerators() ([]User, error) {
	rows, err := db.Query(`SELECT id, username FROM users WHERE role = 'moderator'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moderators []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, err
		}
		moderators = append(moderators, user)
	}
	return moderators, nil
}
