package models

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

	userID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	sessionToken, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	_, err = db.Exec("INSERT INTO users (id, email, username, password, session_token) VALUES (?, ?, ?, ?, ?)",
		userID.String(), email, username, hashedPassword, sessionToken.String())
	if err != nil {
		return "", err
	}
	return sessionToken.String(), err
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

	sessionToken, _ := uuid.NewV4()
	_, err = db.Exec("UPDATE users SET session_token = ? WHERE id = ?", sessionToken.String(), userID)
	if err != nil {
		return "", err
	}

	return sessionToken.String(), nil
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
