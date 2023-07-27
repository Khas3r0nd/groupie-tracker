package internal

import (
	"log"
	"net/http"
	"text/template"
)

func ErrorHandler(w http.ResponseWriter, code int, message string) {
	t, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Println("Error parsing template line 12")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	response := struct {
		ErrorCode int
		ErrorText string
	}{
		ErrorCode: code,
		ErrorText: message,
	}
	t.Execute(w, response)
}
