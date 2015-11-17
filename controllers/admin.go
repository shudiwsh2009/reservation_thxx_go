package controllers

import (
	"encoding/json"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/http"
)

func ViewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	var result = map[string]interface{}{"state": "SUCCESS"}
	var rl = buslogic.ReservationLogic{}

	reservations, err := rl.GetReservationsByAdmin(userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		object := make(map[string]interface{})
		object["reservation_id"] = res.Id
		object["start_time"] = res.StartTime.Format(utils.TIME_PATTERN)
		object["end_time"] = res.EndTime.Format(utils.TIME_PATTERN)
		object["teacher_username"] = res.TeacherUsername
		object["teacher_fullname"] = res.TeacherFullname
		object["teacher_mobile"] = res.TeacherMobile
		object["status"] = res.Status.String()
		array = append(array, object)
	}
	result["reservations"] = array

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}


func ViewMonthlyReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	fromTime := r.PostFormValue("from_time")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var rl = buslogic.ReservationLogic{}

	reservations, err := rl.GetReservationsMonthlyByAdmin(fromTime, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		object := make(map[string]interface{})
		object["reservation_id"] = res.Id
		object["start_time"] = res.StartTime.Format(utils.TIME_PATTERN)
		object["end_time"] = res.EndTime.Format(utils.TIME_PATTERN)
		object["teacher_username"] = res.TeacherUsername
		object["teacher_fullname"] = res.TeacherFullname
		object["teacher_mobile"] = res.TeacherMobile
		object["status"] = res.Status.String()
		array = append(array, object)
	}
	result["reservations"] = array

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}