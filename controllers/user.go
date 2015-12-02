package controllers

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}

	user, err := ul.Login(username, password)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "user_id",
		Value:   fmt.Sprintf("%x", string(user.Id)),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   user.Username,
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "user_type",
		Value:   fmt.Sprintf("%d", user.UserType),
		Path:    "/",
		Expires: time.Now().Local().AddDate(1, 0, 0),
		MaxAge:  365 * 24 * 60 * 60,
	})
	switch user.UserType {
	case models.ADMIN:
		result["url"] = "/appointment/admin"
	case models.TEACHER:
		result["url"] = "/appointment/teacher"
	default:
		result["url"] = "/appointment/entry"
	}

	return result
}

func Logout(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	if userType == models.ADMIN {
		result["url"] = "/appointment/admin"
	} else {
		result["url"] = "/appointment/entry"
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "user_type",
		Path:   "/",
		MaxAge: -1,
	})

	return result
}
