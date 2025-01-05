package handlers

import (
	"forum/internal/models"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

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
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
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
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(24 * time.Hour),
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
