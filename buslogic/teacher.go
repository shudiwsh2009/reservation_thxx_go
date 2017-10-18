package buslogic

import (
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"time"
)

func (w *Workflow) AddReservationByTeacher(startTime string, endTime string, fullname string, mobile string,
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if startTime == "" {
		return nil, re.NewRErrorCodeContext("start_time is empty", nil, re.ErrorMissingParam, "start_time")
	} else if endTime == "" {
		return nil, re.NewRErrorCodeContext("end_time is empty", nil, re.ErrorMissingParam, "end_time")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("start_time is not valid", nil, re.ErrorInvalidParam, "start_time")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("end_time is not valid", nil, re.ErrorInvalidParam, "end_time")
	}
	if start.After(end) {
		return nil, re.NewRErrorCode("start time cannot be after end time", nil, re.ErrorEditReservationEndTimeBeforeStartTime)
	}
	if teacher.Fullname != fullname || teacher.Mobile != mobile {
		teacher.Fullname = fullname
		teacher.Mobile = mobile
		if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to update teacher", err, re.ErrorDatabase)
		}
	}
	// 检查时间是否有冲突
	theDay := utils.BeginOfDay(start)
	nextDay := utils.BeginOfTomorrow(start)
	theDayReservations, err := w.mongoClient.GetReservationsBetweenTime(theDay, nextDay)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get the day reservations", err, re.ErrorDatabase)
	}
	for _, r := range theDayReservations {
		if r.TeacherUsername == teacher.Username {
			if start.After(r.StartTime) && start.Before(r.EndTime) ||
				(end.After(r.StartTime) && end.Before(r.EndTime)) ||
				(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, re.NewRErrorCode("teacher time conflicts", nil, re.ErrorEditReservationTeacherTimeConflict)
			}
		}
	}
	// 新增咨询
	reservation := &model.Reservation{
		StartTime:       start,
		EndTime:         end,
		Status:          model.ReservationStatusAvailable,
		TeacherUsername: teacher.Username,
		TeacherFullname: teacher.Fullname,
		TeacherMobile:   teacher.Mobile,
		TeacherAddress:  teacher.Address,
	}
	if err = w.mongoClient.InsertReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to insert new reservation", err, re.ErrorDatabase)
	}
	return reservation, nil
}

func (w *Workflow) EditReservationByTeacher(reservationId string, startTime string, endTime string,
	fullname string, mobile string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	} else if startTime == "" {
		return nil, re.NewRErrorCodeContext("start_time is empty", nil, re.ErrorMissingParam, "start_time")
	} else if endTime == "" {
		return nil, re.NewRErrorCodeContext("end_time is empty", nil, re.ErrorMissingParam, "end_time")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("cannot get reservation", err, re.ErrorDatabase)
	} else if reservation.Status == model.ReservationStatusReservated {
		return nil, re.NewRErrorCode("cannot edit reservated reservation", nil, re.ErrorEditReservatedReservation)
	} else if reservation.TeacherUsername != teacher.Username {
		return nil, re.NewRErrorCode("cannot edit other's reservation", nil, re.ErrorTeacherEditOtherReservation)
	}
	start, err := time.ParseInLocation("2006-01-02 15:04", startTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("start_time is not valid", nil, re.ErrorInvalidParam, "start_time")
	}
	end, err := time.ParseInLocation("2006-01-02 15:04", endTime, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("end_time is not valid", nil, re.ErrorInvalidParam, "end_time")
	}
	if start.After(end) {
		return nil, re.NewRErrorCode("start time cannot be after end time", nil, re.ErrorEditReservationEndTimeBeforeStartTime)
	} else if start.Before(time.Now()) {
		return nil, re.NewRErrorCode("cannot edit outdated reservation", nil, re.ErrorEditOutdatedReservation)
	}
	if teacher.Fullname != fullname || teacher.Mobile != mobile {
		teacher.Fullname = fullname
		teacher.Mobile = mobile
		if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to update teacher", err, re.ErrorDatabase)
		}
	}
	// 检查时间是否有冲突
	theDay := utils.BeginOfDay(start)
	nextDay := utils.BeginOfTomorrow(start)
	theDayReservations, err := w.mongoClient.GetReservationsBetweenTime(theDay, nextDay)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get the day reservations", err, re.ErrorDatabase)
	}
	for _, r := range theDayReservations {
		if r.TeacherUsername == teacher.Username && r.Id.Hex() != reservation.Id.Hex() {
			if start.After(r.StartTime) && start.Before(r.EndTime) ||
				(end.After(r.StartTime) && end.Before(r.EndTime)) ||
				(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, re.NewRErrorCode("teacher time conflicts", nil, re.ErrorEditReservationTeacherTimeConflict)
			}
		}
	}
	// 更新咨询
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherFullname = teacher.Fullname
	reservation.TeacherMobile = teacher.Mobile
	if err = w.mongoClient.UpdateReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation", err, re.ErrorDatabase)
	}
	return reservation, nil
}

func (w *Workflow) RemoveReservationsByTeacher(reservationIds []string, userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return 0, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return 0, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	removed := 0
	for _, reservationId := range reservationIds {
		reservation, err := w.mongoClient.GetReservationById(reservationId)
		if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted ||
			reservation.TeacherUsername != teacher.Username {
			continue
		}
		reservation.Status = model.ReservationStatusDeleted
		if w.mongoClient.UpdateReservation(reservation) == nil {
			removed++
		}
	}
	return removed, nil
}

func (w *Workflow) CancelReservationsByTeacher(reservationIds []string, userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return 0, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return 0, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	canceled := 0
	for _, reservationId := range reservationIds {
		reservation, err := w.mongoClient.GetReservationById(reservationId)
		if err != nil || reservation == nil || reservation.Status != model.ReservationStatusReservated ||
			reservation.StartTime.Before(time.Now()) || reservation.TeacherUsername != teacher.Username {
			continue
		}
		reservation.Status = model.ReservationStatusAvailable
		reservation.StudentInfo = model.StudentInfo{}
		reservation.StudentFeedback = model.StudentFeedback{}
		reservation.TeacherFeedback = model.TeacherFeedback{}
		if w.mongoClient.UpdateReservation(reservation) == nil {
			canceled++
		}
	}
	return canceled, nil
}

func (w *Workflow) GetFeedbackByTeacher(reservationId string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ErrorFeedbackFutureReservation)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorFeedbackAvailableReservation)
	} else if reservation.TeacherUsername != teacher.Username {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ErrorFeedbackOtherReservation)
	}
	return reservation, nil
}

func (w *Workflow) SubmitFeedbackByTeacher(reservationId string, teacherFullname string, teacherUsername string,
	studentFullname string, problem string, solution string, adviceToCenter string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	} else if teacherFullname == "" {
		return nil, re.NewRErrorCodeContext("teacher_fullname is empty", nil, re.ErrorMissingParam, "teacher_fullname")
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher_username is empty", nil, re.ErrorMissingParam, "teacher_username")
	} else if studentFullname == "" {
		return nil, re.NewRErrorCodeContext("student_fullname is empty", nil, re.ErrorMissingParam, "student_fullname")
	} else if problem == "" {
		return nil, re.NewRErrorCodeContext("problem is empty", nil, re.ErrorMissingParam, "problem")
	} else if solution == "" {
		return nil, re.NewRErrorCodeContext("solution is empty", nil, re.ErrorMissingParam, "solution")
	} else if adviceToCenter == "" {
		return nil, re.NewRErrorCodeContext("advice_to_center is empty", nil, re.ErrorMissingParam, "advice_to_center")
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ErrorFeedbackFutureReservation)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorFeedbackAvailableReservation)
	} else if reservation.TeacherUsername != teacher.Username {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ErrorFeedbackOtherReservation)
	}
	sendFeedbackSms := reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty()
	reservation.TeacherFeedback = model.TeacherFeedback{
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherUsername,
		StudentFullname: studentFullname,
		Problem:         problem,
		Solution:        solution,
		AdviceToCenter:  adviceToCenter,
	}
	if err = w.mongoClient.UpdateReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation", err, re.ErrorDatabase)
	}
	if sendFeedbackSms {
		go w.SendFeedbackSMS(reservation)
	}
	return reservation, nil
}

func (w *Workflow) GetReservationStudentInfoByTeacher(reservationId string, userId string, userType int) (*model.StudentInfo, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorViewAvailableReservationStudentInfo)
	} else if reservation.TeacherUsername != teacher.Username {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ErrorTeacherViewOtherReservation)
	}
	return &reservation.StudentInfo, nil
}

func (w *Workflow) WrapSimpleTeacher(teacher *model.Teacher) map[string]interface{} {
	var result = make(map[string]interface{})
	if teacher == nil {
		return result
	}
	result["fullname"] = teacher.Fullname
	result["address"] = teacher.Address
	result["gender"] = teacher.Gender
	result["major"] = teacher.Major
	result["academic"] = teacher.Academic
	result["aptitude"] = teacher.Aptitude
	result["problem"] = teacher.Problem
	return result
}

func (w *Workflow) WrapTeacher(teacher *model.Teacher) map[string]interface{} {
	var result = w.WrapSimpleTeacher(teacher)
	if teacher == nil {
		return result
	}
	result["username"] = teacher.Username
	result["mobile"] = teacher.Mobile
	return result
}
