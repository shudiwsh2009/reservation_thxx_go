package buslogic

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"time"
)

func (w *Workflow) MakeReservationByStudent(reservationId string, fullname string, gender string,
	username string, school string, hometown string, mobile string, email string, experience string,
	problem string) (*model.Reservation, error) {
	if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ErrorMissingParam, "reservation_id")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if gender == "" {
		return nil, re.NewRErrorCodeContext("gender is empty", nil, re.ErrorMissingParam, "gender")
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if school == "" {
		return nil, re.NewRErrorCodeContext("school is empty", nil, re.ErrorMissingParam, "school")
	} else if hometown == "" {
		return nil, re.NewRErrorCodeContext("hometown is empty", nil, re.ErrorMissingParam, "hometown")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if email == "" {
		return nil, re.NewRErrorCodeContext("email is empty", nil, re.ErrorMissingParam, "email")
	} else if experience == "" {
		return nil, re.NewRErrorCodeContext("experience is empty", nil, re.ErrorMissingParam, "experience")
	} else if problem == "" {
		return nil, re.NewRErrorCodeContext("problem is empty", nil, re.ErrorMissingParam, "problem")
	} else if !utils.IsStudentUsername(username) {
		return nil, re.NewRErrorCode("student username format is wrong", nil, re.ErrorFormatStudentUsername)
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	} else if !utils.IsEmail(email) {
		return nil, re.NewRErrorCode("email format is wrong", nil, re.ErrorFormatEmail)
	} else if config.Instance().StudentVerificationEnabled && !w.verifyStudent(username, fullname) {
		return nil, re.NewRErrorCode("student verification failed", nil, re.ErrorStudentFullnameNotMatch)
	}

	studentReservations, err := w.MongoClient().GetReservationsByStudentUsername(username)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ErrorDatabase)
	}
	for _, r := range studentReservations {
		if r.Status == model.ReservationStatusReservated && r.StartTime.After(time.Now()) {
			return nil, re.NewRErrorCode("already have reservation", nil, re.ErrorStudentAlreadyHaveReservation)
		}
	}

	reservation, err := w.MongoClient().GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", nil, re.ErrorDatabase)
	} else if reservation.StartTime.Before(time.Now()) {
		return nil, re.NewRErrorCode("cannot make outdated reservation", nil, re.ErrorStudentMakeOutdatedReservation)
	} else if reservation.Status != model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot make reservated reservation", nil, re.ErrorStudentMakeReservatedReservation)
	} else if time.Now().Add(model.MakeReservationLatestDuration).After(reservation.StartTime) {
		return nil, re.NewRErrorCodeContext("cannot make reservation starting in 3 hours", nil,
			re.ErrorStudentMakeReservationTooEarly, fmt.Sprintf("%d小时", int64(model.MakeReservationLatestDuration.Hours())))
	}

	reservation.Status = model.ReservationStatusReservated
	reservation.StudentInfo = model.StudentInfo{
		Fullname:   fullname,
		Gender:     gender,
		Username:   username,
		School:     school,
		Hometown:   hometown,
		Mobile:     mobile,
		Email:      email,
		Experience: experience,
		Problem:    problem,
	}
	err = w.MongoClient().UpdateReservation(reservation)
	if err != nil {
		return nil, re.NewRErrorCode("fail to update reservation", err, re.ErrorDatabase)
	}

	teacher, err := w.MongoClient().GetTeacherByUsername(reservation.TeacherUsername)
	if err != nil || teacher == nil {
		return nil, re.NewRErrorCode("failed to GetTeacherByUsername", nil, re.ErrorDatabase)
	}

	//send success sms
	go w.SendSuccessSMS(reservation, teacher)
	return reservation, nil
}

func (w *Workflow) GetFeedbackByStudent(reservationId string, username string) (*model.Reservation, error) {
	if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ErrorMissingParam, "reservation_id")
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if !utils.IsStudentUsername(username) {
		return nil, re.NewRErrorCode("student username format is wrong", nil, re.ErrorFormatStudentUsername)
	}

	reservation, err := w.MongoClient().GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ErrorFeedbackFutureReservation)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorFeedbackAvailableReservation)
	} else if reservation.StudentInfo.Username != username {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ErrorFeedbackOtherReservation)
	}
	return reservation, nil
}

func (w *Workflow) SubmitFeedbackByStudent(reservationId string, fullname string, problem string, choices string,
	score string, feedback string, username string) (*model.Reservation, error) {
	if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ErrorMissingParam, "reservation_id")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if problem == "" {
		return nil, re.NewRErrorCodeContext("problem is empty", nil, re.ErrorMissingParam, "problem")
	} else if choices == "" {
		return nil, re.NewRErrorCodeContext("choices is empty", nil, re.ErrorMissingParam, "choices")
	} else if score == "" {
		return nil, re.NewRErrorCodeContext("score is empty", nil, re.ErrorMissingParam, "score")
	} else if feedback == "" {
		return nil, re.NewRErrorCodeContext("feedback is empty", nil, re.ErrorMissingParam, "feedback")
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if !utils.IsStudentUsername(username) {
		return nil, re.NewRErrorCode("student username format is wrong", nil, re.ErrorFormatStudentUsername)
	}

	reservation, err := w.MongoClient().GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ErrorFeedbackFutureReservation)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorFeedbackAvailableReservation)
	} else if reservation.StudentInfo.Username != username {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ErrorFeedbackOtherReservation)
	}

	reservation.StudentFeedback = model.StudentFeedback{
		Fullname: fullname,
		Problem:  problem,
		Choices:  choices,
		Score:    score,
		Feedback: feedback,
	}
	if err = w.MongoClient().UpdateReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation", err, re.ErrorDatabase)
	}
	return reservation, nil
}

func (w *Workflow) GetReservationTeacherInfoByStudent(reservationId string) (*model.Teacher, error) {
	if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation id is empty", nil, re.ErrorMissingParam, "reservation_id")
	}

	reservation, err := w.MongoClient().GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	}
	teacher, err := w.MongoClient().GetTeacherByUsername(reservation.TeacherUsername)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	return teacher, nil
}

func (w *Workflow) WrapStudenInfo(studentInfo *model.StudentInfo) map[string]interface{} {
	var result = make(map[string]interface{})
	if studentInfo == nil {
		return result
	}
	result["fullname"] = studentInfo.Fullname
	result["gender"] = studentInfo.Gender
	result["username"] = studentInfo.Username
	result["school"] = studentInfo.School
	result["hometown"] = studentInfo.Hometown
	result["mobile"] = studentInfo.Mobile
	result["email"] = studentInfo.Email
	result["experience"] = studentInfo.Experience
	result["problem"] = studentInfo.Problem
	return result
}
