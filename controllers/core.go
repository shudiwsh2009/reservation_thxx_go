package controllers

import (
	"encoding/json"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"html/template"
	"net/http"
)

func EntryPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/entry.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func StudentPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/student.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func LoginPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) interface{} {
	t := template.Must(template.ParseFiles("../templates/login.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func TeacherPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) interface{} {
	if userType == models.ADMIN {
		http.Redirect(w, r, "/appointment/admin", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("../templates/teacher.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func AdminPage(w http.ResponseWriter, r *http.Request, username string, userType models.UserType) interface{} {
	if userType == models.TEACHER {
		http.Redirect(w, r, "/appointment/teacher", http.StatusFound)
		return nil
	}
	t := template.Must(template.ParseFiles("../templates/admin.html"))
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
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
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}
