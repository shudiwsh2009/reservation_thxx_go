package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func EntryPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/entry.html"))
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {

	}
}

func StudentPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/student.html"))
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {

	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/login.html"))
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {

	}
}

func TeacherPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/teacher.html"))
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {

	}
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/admin.html"))
		err := t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {

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
