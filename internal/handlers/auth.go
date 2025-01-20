package handlers

import (
	"context"
	"encoding/json"
	"forum/internal/models"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

var googleOAuthConfig = &oauth2.Config{
	ClientID:     "9254237354-63n6pgm41osqflqoq98boibdnkpsuggv.apps.googleusercontent.com",
	ClientSecret: "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",
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Google OAuth2 exchange failed:", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	client := googleOAuthConfig.Client(context.Background(), token)
	userInfo, err := getGoogleUserInfo(client)
	if err != nil {
		log.Println("Failed to get Google user info:", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Проверка или создание пользователя в базе данных
	sessionToken, err := models.AuthenticateOrRegisterOAuthUser(userInfo.Email, userInfo.Name, "google")
	if err != nil {
		log.Println("Failed to authenticate/register user:", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// Установка куки с токеном сессии
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",                  // Доступ к куки на всех маршрутах
		HttpOnly: true,                 // Защита от доступа через JavaScript
		Secure:   false,                // Для локального тестирования
		SameSite: http.SameSiteLaxMode, // Уменьшение риска CSRF
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getGoogleUserInfo(client *http.Client) (*models.OAuthUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo models.OAuthUserInfo
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	return &userInfo, err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := sanitizeInput(r.FormValue("email"))
		username := sanitizeInput(r.FormValue("username"))
		password := r.FormValue("password")

		if email == "" || username == "" || password == "" {
			ErrorHandler(w, r, http.StatusBadRequest, "All fields are required")
			return
		}

		if !isValidEmail(email) {
			ErrorHandler(w, r, http.StatusBadRequest, "Invalid email format")
			return
		}

		emailExists, err := models.CheckEmailExists(email)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if emailExists {
			renderTemplateWithError(w, "register", "Email is already registered")
			return
		}

		usernameExists, err := models.CheckUsernameExists(username)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		if usernameExists {
			renderTemplateWithError(w, "register", "Username is already taken")
			return
		}

		sessionToken, err := models.RegisterUser(email, username, password)
		if err != nil {
			ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		cookie := http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",                  // Доступ к куки на всех маршрутах
			HttpOnly: true,                 // Защита от доступа через JavaScript
			Secure:   false,                // Для локального тестирования
			SameSite: http.SameSiteLaxMode, // Уменьшение риска CSRF
		}

		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/register.html")
	tmpl.Execute(w, nil)
}

func sanitizeInput(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return r
		}
		return -1
	}, trimmed)
}

func renderTemplateWithError(w http.ResponseWriter, templateName, errorMessage string) {
	tmpl, _ := template.ParseFiles("web/templates/" + templateName + ".html")
	tmpl.Execute(w, struct{ Error string }{Error: errorMessage})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		sessionToken, err := models.AuthenticateUser(email, password)
		if err != nil {
			tmpl, _ := template.ParseFiles("web/templates/login.html")
			tmpl.Execute(w, struct{ Error string }{Error: "Invalid email or password"})
			return
		}

		cookie := http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",                  // Доступ к куки на всех маршрутах
			HttpOnly: true,                 // Защита от доступа через JavaScript
			Secure:   false,                // Для локального тестирования
			SameSite: http.SameSiteLaxMode, // Уменьшение риска CSRF
		}

		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/login.html")
	tmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}
