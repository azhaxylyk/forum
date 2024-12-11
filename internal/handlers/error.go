package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

type ErrorPage struct {
	StatusCode string
	StatusText string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int, statusText string) {
	data := ErrorPage{
		StatusCode: fmt.Sprint(statusCode),
		StatusText: statusText,
	}

	ts, err := template.ParseFiles("web/templates/wentwrong.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		http.Error(w, "Error when executing", http.StatusInternalServerError)
		return
	}
}
