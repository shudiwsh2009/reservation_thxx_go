package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/http"
	"strconv"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}

	user, err := ul.Login(username, password)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   fmt.Sprintf("%x", string(user.Id)),
		Path: "/",
		Expires: time.Now().In(utils.Location).AddDate(1, 0, 0),
		MaxAge: 365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   user.Username,
		Path: "/",
		Expires: time.Now().In(utils.Location).AddDate(1, 0, 0),
		MaxAge: 365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_type",
		Value:   fmt.Sprintf("%d", user.UserType),
		Path: "/",
		Expires: time.Now().In(utils.Location).AddDate(1, 0, 0),
		MaxAge: 365 * 24 * 60 * 60,
	})
	switch user.UserType {
	case models.ADMIN:
		result["url"] = "/appointment/admin"
	case models.TEACHER:
		result["url"] = "/appointment/teacher"
	default:
		result["url"] = "/appointment/entry"
	}

	if data, err := json.Marshal(result); err == nil {
		fmt.Println(string(data))
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var result = map[string]interface{}{"state": "SUCCESS"}

	userTypeCookie, _ := r.Cookie("user_type")
	if userType, err := strconv.Atoi(userTypeCookie.Value); err != nil && models.UserType(userType) == models.ADMIN {
		result["url"] = "/appointment/admin"
	} else {
		result["url"] = "/appointment/entry"
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Path: "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Path: "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_type",
		Path: "/",
		MaxAge: -1,
	})

	if data, err := json.Marshal(result); err == nil {
		fmt.Println(string(data))
		w.Write(data)
	}
}
