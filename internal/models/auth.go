package models

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type OAuthUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func AuthenticateOrRegisterOAuthUser(email, username, provider string) (string, error) {
	var userID string

	// Проверяем, существует ли пользователь с таким email
	err := db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователя нет, регистрируем его
			newUUID, err := uuid.NewV4() // Создаем UUID
			if err != nil {
				return "", err
			}
			userID = newUUID.String() // Преобразуем UUID в строку

			sessionUUID, err := uuid.NewV4() // Создаем токен сессии
			if err != nil {
				return "", err
			}
			sessionToken := sessionUUID.String() // Преобразуем UUID в строку

			_, err = db.Exec("INSERT INTO users (id, email, username, provider, session_token) VALUES (?, ?, ?, ?, ?)",
				userID, email, username, provider, sessionToken)
			if err != nil {
				return "", err
			}
			return sessionToken, nil
		}
		return "", err
	}

	// Если пользователь уже существует, обновляем токен сессии
	sessionUUID, err := uuid.NewV4() // Создаем новый токен сессии
	if err != nil {
		return "", err
	}
	sessionToken := sessionUUID.String() // Преобразуем UUID в строку

	_, err = db.Exec("UPDATE users SET session_token = ? WHERE id = ?", sessionToken, userID)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	return exists, err
}

func CheckUsernameExists(username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	return exists, err
}

func RegisterUser(email, username, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	userID := newUUID.String()

	sessionUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionToken := sessionUUID.String()

	_, err = db.Exec("INSERT INTO users (id, email, username, password, session_token) VALUES (?, ?, ?, ?, ?)",
		userID, email, username, hashedPassword, sessionToken)
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func AuthenticateUser(email, password string) (string, error) {
	var userID, hashedPassword string

	err := db.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	sessionUUID, _ := uuid.NewV4()
	sessionToken := sessionUUID.String()

	_, err = db.Exec("UPDATE users SET session_token = ? WHERE id = ?", sessionToken, userID)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func GetIDBySessionToken(sessionToken string) (string, string, error) {
	var username string
	var userID string
	err := db.QueryRow("SELECT id, username FROM users WHERE session_token = ?", sessionToken).Scan(&userID, &username)
	if err != nil {
		return "", "", err
	}
	return userID, username, nil
}
