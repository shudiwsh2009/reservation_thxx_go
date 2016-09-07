package controllers

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/http"
	"net/url"
	"time"
)

func ViewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}
	var rl = buslogic.ReservationLogic{}

	reservations, err := rl.GetReservationsByAdmin(userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["end_time"] = res.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["teacher_username"] = res.TeacherUsername
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

	return result
}

func ViewMonthlyReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["from_time"]) == 0 {
		ErrorHandler(w, r, errors.New("参数错误"))
		return nil
	}
	fromTime := queryForm["from_time"][0]

	var result = map[string]interface{}{"state": "SUCCESS"}
	var rl = buslogic.ReservationLogic{}

	reservations, err := rl.GetReservationsMonthlyByAdmin(fromTime, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id
		resJson["start_time"] = res.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["end_time"] = res.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
		resJson["teacher_username"] = res.TeacherUsername
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

	return result
}

func AddReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var reservationJson = make(map[string]interface{})
	reservation, err := al.AddReservationByAdmin(startTime, endTime, teacherUsername, teacherFullname,
		teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["teacher_username"] = reservation.TeacherUsername
	reservationJson["teacher_fullname"] = reservation.TeacherFullname
	reservationJson["teacher_mobile"] = reservation.TeacherMobile
	result["reservation"] = reservationJson

	return result
}

func EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	startTime := r.PostFormValue("start_time")
	endTime := r.PostFormValue("end_time")
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMobile := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var reservationJson = make(map[string]interface{})
	reservation, err := al.EditReservationByAdmin(reservationId, startTime, endTime, teacherUsername,
		teacherFullname, teacherMobile, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	reservationJson["reservation_id"] = reservation.Id
	reservationJson["start_time"] = reservation.StartTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["end_time"] = reservation.EndTime.In(utils.Location).Format(utils.TIME_PATTERN)
	reservationJson["teacher_username"] = reservation.TeacherUsername
	reservationJson["teacher_fullname"] = reservation.TeacherFullname
	reservationJson["teacher_mobile"] = reservation.TeacherMobile
	result["reservation"] = reservationJson

	return result
}

func RemoveReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	removed, err := al.RemoveReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}

func CancelReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	removed, err := al.CancelReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["removed_count"] = removed

	return result
}

func GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var feedback = make(map[string]interface{})
	reservation, err := al.GetFeedbackByAdmin(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
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

	return result
}

func SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
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
		return nil
	}

	return result
}

func GetStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var studentJson = make(map[string]interface{})
	studentInfo, err := al.GetStudentInfoByAdmin(reservationId, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
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

	return result
}

func ExportReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	r.ParseForm()
	reservationIds := []string(r.Form["reservation_ids"])

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	url, err := al.ExportReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	result["url"] = url

	return result
}

func SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	teacherUsername := r.PostFormValue("teacher_username")
	teacherFullname := r.PostFormValue("teacher_fullname")
	teacherMoble := r.PostFormValue("teacher_mobile")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var teacherJson = make(map[string]interface{})
	teacher, err := al.SearchTeacherByAdmin(teacherFullname, teacherUsername, teacherMoble,
		userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	teacherJson["teacher_username"] = teacher.Username
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher"] = teacherJson

	return result
}

func GetTeacherInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	teacherUsername := r.PostFormValue("teacher_username")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	var teacherJson = make(map[string]interface{})
	teacher, err := al.GetTeacherInfoByAdmin(teacherUsername, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}
	teacherJson["teacher_username"] = teacher.Username
	teacherJson["fullname"] = teacher.Fullname
	teacherJson["gender"] = teacher.Gender
	teacherJson["major"] = teacher.Major
	teacherJson["academic"] = teacher.Academic
	teacherJson["aptitude"] = teacher.Aptitude
	teacherJson["problem"] = teacher.Problem
	result["teacher"] = teacherJson

	return result
}

func EditTeacherInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType models.UserType) interface{} {
	teacherUsername := r.PostFormValue("teacher_username")
	fullname := r.PostFormValue("fullname")
	gender := r.PostFormValue("gender")
	major := r.PostFormValue("major")
	academic := r.PostFormValue("academic")
	aptitude := r.PostFormValue("aptitude")
	problem := r.PostFormValue("problem")

	var result = map[string]interface{}{"state": "SUCCESS"}
	var al = buslogic.AdminLogic{}

	_, err := al.EditTeacherInfoByAdmin(teacherUsername, fullname, gender, major, academic, aptitude, problem, userId, userType)
	if err != nil {
		ErrorHandler(w, r, err)
		return nil
	}

	return result
}
