package models

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type Notification struct {
	ID                 string
	UserID             string
	ActionBy           string
	ActionType         string
	TargetID           string
	TargetType         string
	IsRead             bool
	CreatedAt          time.Time
	CreatedAtFormatted string
}

// CreateNotification создает новое уведомление
func CreateNotification(userID, actionBy, actionType, targetID, targetType string) error {
	notificationID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO notifications (id, user_id, action_by, action_type, target_id, target_type, is_read, created_at)
		VALUES (?, ?, ?, ?, ?, ?, FALSE, ?)
	`, notificationID.String(), userID, actionBy, actionType, targetID, targetType, time.Now())
	return err
}

// GetNotificationsForUser возвращает уведомления для пользователя
func GetNotificationsForUser(userID string) ([]Notification, error) {
	rows, err := db.Query(`
		SELECT id, action_by, action_type, target_id, target_type, is_read, created_at
		FROM notifications
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		fmt.Printf("Error querying notifications for user %s: %v\n", userID, err)
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		var createdAt time.Time
		err := rows.Scan(
			&notification.ID,
			&notification.ActionBy,
			&notification.ActionType,
			&notification.TargetID,
			&notification.TargetType,
			&notification.IsRead,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		notification.CreatedAt = createdAt
		notification.CreatedAtFormatted = createdAt.Format("02.01.2006 15:04")
		fmt.Printf("Fetched notification: %+v\n", notification)
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

// MarkNotificationAsRead помечает уведомление как прочитанное
func MarkNotificationAsRead(notificationID string) error {
	_, err := db.Exec(`
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = ?
	`, notificationID)
	return err
}

// MarkAllNotificationsAsRead помечает все уведомления пользователя как прочитанные
func MarkAllNotificationsAsRead(userID string) error {
	_, err := db.Exec(`
		UPDATE notifications
		SET is_read = TRUE
		WHERE user_id = ?
	`, userID)
	return err
}

// DeleteReadNotifications удаляет прочитанные уведомления пользователя
func DeleteReadNotifications(userID string) error {
	_, err := db.Exec(`
		DELETE FROM notifications
		WHERE user_id = ? AND is_read = TRUE
	`, userID)
	return err
}

func (n *Notification) GetMessage() (string, error) {
	// Получаем имя пользователя, который совершил действие
	username, err := GetUsernameByID(n.ActionBy)
	if err != nil {
		return "", err
	}

	// Формируем сообщение на основе типа действия и цели
	switch n.ActionType {
	case "comment":
		return fmt.Sprintf("%s commented on your post.", username), nil
	case "like":
		if n.TargetType == "post" {
			return fmt.Sprintf("%s liked your post.", username), nil
		} else if n.TargetType == "comment" {
			return fmt.Sprintf("%s liked your comment.", username), nil
		}
	case "dislike":
		if n.TargetType == "post" {
			return fmt.Sprintf("%s disliked your post.", username), nil
		} else if n.TargetType == "comment" {
			return fmt.Sprintf("%s disliked your comment.", username), nil
		}
	case "approve_del":
		// Сообщение о том, что пост был удален администратором
		return "Your request has been approved", nil
	case "approve_mod":
		// Сообщение о том, что пользователь был утвержден как модератор
		return "Admin has approved you as a moderator.", nil
	case "reject_mod":
		// Сообщение о том, что пользователь был отклонен как модератор
		return "Admin has rejected your application to be a moderator.", nil
	case "reject_del":
		// Сообщение о том, что запрос на удаление поста был отклонен
		return "Your request to delete the post has been rejected.", nil
	}

	return "You have a new notification.", nil
}

// GetUsernameByID получает имя пользователя по его ID
func GetUsernameByID(userID string) (string, error) {
	var username string
	err := db.QueryRow(`SELECT username FROM users WHERE id = ?`, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

// GetUnreadNotificationCount возвращает количество непрочитанных уведомлений для пользователя
func GetUnreadNotificationCount(userID string) (int, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM notifications 
		WHERE user_id = ? AND is_read = FALSE
	`, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
