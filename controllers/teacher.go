package controllers

import (
	"encoding/json"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/http"
	"time"
)

func ViewReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	var result = map[string]interface{}{"state": "SUCCESS"}
	var ul = buslogic.UserLogic{}
	var rl = buslogic.ReservationLogic{}

	teacher, err := ul.GetUserById(userId)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	var teacherJson = make(map[string]interface{})
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher_info"] = teacherJson

	reservations, err := rl.GetReservationsByTeacher(userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["end_time"] = res.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["teacher_fullname"] = res.TeacherFullname
		resJson["teacher_mobile"] = res.TeacherMobile
		if res.Status == models.AVAILABLE {
			resJson["status"] = models.AVAILABLE.String()
		} else if res.Status == models.RESERVATED && res.StartTime.Before(time.Now().In(utils.Location)) {
			resJson["status"] = models.FEEDBACK.String()
		} else {
			resJson["status"] = models.RESERVATED.String()
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func AddReservationByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	var reservationJson = make(map[string]interface{})
	reservation, err := tl.AddReservationByTeacher(startTime, endTime, teacherFullname, teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["teacher_fullname"] = reservation.TeacherFullname
	reservationJson["teacher_mobile"] = reservation.TeacherMobile
	result["reservation"] = reservationJson

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func EditReservationByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	var reservationJson = make(map[string]interface{})
	reservation, err := tl.EditReservationByTeacher(reservationId, startTime, endTime, teacherFullname, teacherMobile,
		userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["teacher_fullname"] = reservation.TeacherFullname
	reservationJson["teacher_mobile"] = reservation.TeacherMobile
	result["reservation"] = reservationJson

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func RemoveReservationByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	removed, err := tl.RemoveReservationsByTeacher(reservationIds, userId, userType)
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

func CancelReservationByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	removed, err := tl.CancelReservationsByTeacher(reservationIds, userId, userType)
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

func GetFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	var feedback = make(map[string]interface{})
	reservation, err := tl.GetFeedbackByTeacher(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	if len(reservation.TeacherFeedback.TeacherFullname) == 0 {
		feedback["teacher_fullname"] = reservation.TeacherFullname
	} else {
		feedback["teacher_fullname"] = reservation.TeacherFeedback.TeacherFullname
	}
	if len(reservation.TeacherFeedback.TeacherUsername) == 0 {
		feedback["teacher_username"] = reservation.TeacherUsername
	} else {
		feedback["teacher_username"] = reservation.TeacherFeedback.TeacherUsername
	}
	if len(reservation.TeacherFeedback.StudentFullname) == 0 {
		feedback["student_fullname"] = reservation.StudentInfo.Name
	} else {
		feedback["student_fullname"] = reservation.TeacherFeedback.StudentFullname
	}
	feedback["problem"] = reservation.TeacherFeedback.Problem
	feedback["solution"] = reservation.TeacherFeedback.Solution
	feedback["advice"] = reservation.TeacherFeedback.AdviceToCenter
	result["feedback"] = feedback

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func SubmitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	studentFullname := r.PostFormValue("student_fullname")
	problem := r.PostFormValue("problem")
	solution := r.PostFormValue("solution")
	advice := r.PostFormValue("advice")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	_, err := tl.SubmitFeedbackByTeacher(reservationId, teacherFullname, teacherUsername, studentFullname,
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

func GetStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var tl = buslogic.TeacherLogic{}

	var studentJson = make(map[string]interface{})
	studentInfo, err := tl.GetStudentInfoByTeacher(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	studentJson["name"] = studentInfo.Name
	studentJson["gender"] = studentInfo.Gender
	studentJson["student_id"] = studentInfo.StudentId
	studentJson["school"] = studentInfo.School
	studentJson["hometown"] = studentInfo.Hometown
	studentJson["mobile"] = studentInfo.Mobile
	studentJson["email"] = studentInfo.Email
	studentJson["experience"] = studentInfo.Experience
	studentJson["problem"] = studentInfo.Problem
	result["student_info"] = studentJson

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}
