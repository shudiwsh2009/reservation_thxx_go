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
		object["start_time"] = res.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
		object["end_time"] = res.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
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
		object["start_time"] = res.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
		object["end_time"] = res.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
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

func AddReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	reservation, err := al.AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	result["reservation_id"] = reservation.Id
	result["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	result["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	result["teacher_username"] = reservation.TeacherUsername
	result["teacher_fullname"] = reservation.TeacherFullname
	result["teacher_mobile"] = reservation.TeacherMobile

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	reservation, err := al.EditReservationByAdmin(reservationId, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	result["reservation_id"] = reservation.Id
	result["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	result["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	result["teacher_username"] = reservation.TeacherUsername
	result["teacher_fullname"] = reservation.TeacherFullname
	result["teacher_mobile"] = reservation.TeacherMobile

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func RemoveReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	removed, err := al.RemoveReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	result["removed_count"] = removed

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func CancelReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	removed, err := al.CancelReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	result["removed_count"] = removed

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	reservation, err := al.GetFeedbackByAdmin(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	if len(reservation.TeacherFeedback.TeacherFullname) == 0 {
		result["teacher_fullname"] = reservation.TeacherFullname
	} else {
		result["teacher_fullname"] = reservation.TeacherFeedback.TeacherFullname
	}
	if len(reservation.TeacherFeedback.TeacherUsername) == 0 {
		result["teacher_username"] = reservation.TeacherUsername
	} else {
		result["teacher_username"] = reservation.TeacherFeedback.TeacherUsername
	}
	if len(reservation.TeacherFeedback.StudentFullname) == 0 {
		result["student_fullname"] = reservation.StudentInfo.Name
	} else {
		result["student_fullname"] = reservation.TeacherFeedback.StudentFullname
	}
	result["problem"] = reservation.TeacherFeedback.Problem
	result["solution"] = reservation.TeacherFeedback.Solution
	result["advice"] = reservation.TeacherFeedback.AdviceToCenter

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	studentFullname := r.PostFormValue("student_fullname")
	problem := r.PostFormValue("problem")
	solution := r.PostFormValue("solution")
	advice := r.PostFormValue("advice")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	_, err := al.SubmitFeedbackByAdmin(reservationId, teacherFullname, teacherUsername, studentFullname,
		problem, solution, advice, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMoble := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	teacher, err := al.SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	result["teacher_username"] = teacher.Username
	result["teacher_fullname"] = teacher.Fullname
	result["teacher_mobile"] = teacher.Mobile

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}