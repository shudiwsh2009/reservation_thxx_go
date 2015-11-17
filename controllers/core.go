package controllers

import (
	"encoding/json"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"html/template"
	"net/http"
)

func EntryPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) {
	t := template.Must(template.ParseFiles("templates/entry.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StudentPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) {
	t := template.Must(template.ParseFiles("templates/student.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) {
	t := template.Must(template.ParseFiles("templates/login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TeacherPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) {
	t := template.Must(template.ParseFiles("templates/teacher.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AdminPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) {
	t := template.Must(template.ParseFiles("templates/admin.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type ErrorMsg struct {
	State   string `json:"state"`
	Message string `json:"message"`
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if data, err := json.Marshal(ErrorMsg{
		State:   "FAILED",
		Message: err.Error(),
	}); err == nil {
		w.Write(data)
	}
}
