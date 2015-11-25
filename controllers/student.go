package controllers

import (
	"encoding/json"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/http"
	"time"
)

func ViewReservationsByStudent(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	var result = map[string]interface{}{"state": "SUCCESS"}
	var rl = buslogic.ReservationLogic{}

	reservations, err := rl.GetReservationsByStudent()
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

func MakeReservationByStudent(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")
	name := r.PostFormValue("name")
	gender := r.PostFormValue("gender")
	studentId := r.PostFormValue("student_id")
	school := r.PostFormValue("school")
	hometown := r.PostFormValue("hometown")
	mobile := r.PostFormValue("mobile")
	email := r.PostFormValue("email")
	experience := r.PostFormValue("experience")
	problem := r.PostFormValue("problem")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var sl = buslogic.StudentLogic{}

	var reservationJson = make(map[string]interface{})
	reservation, err := sl.MakeReservationByStudent(reservationId, name, gender, studentId, school, hometown,
		mobile, email, experience, problem)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	reservationJson["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["teacher_fullname"] = reservation.TeacherFullname
	result["reservation"] = reservationJson

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func GetFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var sl = buslogic.StudentLogic{}

	var feedbackJson = make(map[string]interface{})
	reservation, err := sl.GetFeedbackByStudent(reservationId, studentId)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}
	if len(reservation.StudentFeedback.Name) == 0 {
		feedbackJson["name"] = reservation.StudentInfo.Name
	} else {
		feedbackJson["name"] = reservation.StudentFeedback.Name
	}
	if len(reservation.StudentFeedback.Problem) == 0 {
		feedbackJson["problem"] = reservation.StudentInfo.Problem
	} else {
		feedbackJson["problem"] = reservation.StudentFeedback.Problem
	}
	feedbackJson["choices"] = reservation.StudentFeedback.Choices
	feedbackJson["score"] = reservation.StudentFeedback.Score
	feedbackJson["feedback"] = reservation.StudentFeedback.Feedback
	result["feedback"] = feedbackJson

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}

func SubmitFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) {
	reservationId := r.PostFormValue("reservation_id")
	name := r.PostFormValue("name")
	problem := r.PostFormValue("problem")
	choices := r.PostFormValue("choices")
	score := r.PostFormValue("score")
	feedback := r.PostFormValue("feedback")
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var sl = buslogic.StudentLogic{}

	_, err := sl.SubmitFeedbackByStudent(reservationId, name, problem, choices, score, feedback, studentId)
	if err != nil {
		ErrorHandler(w, r, err)
		return
	}

	if data, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Write(data)
	}
}
