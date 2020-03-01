package buslogic

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"time"
)

func (w *Workflow) AddReservationByTeacher(startTime string, endTime string, fullname string, fullnameEn string,
	mobile string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if startTime == "" {
		return nil, re.NewRErrorCodeContext("start_time is empty", nil, re.ErrorMissingParam, "start_time")
	} else if endTime == "" {
		return nil, re.NewRErrorCodeContext("end_time is empty", nil, re.ErrorMissingParam, "end_time")
	} else if fullname == "" && fullnameEn == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	}
	teacher, err := w.MongoClient().GetTeacherById(userId)
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
	if teacher.Fullname != fullname || teacher.FullnameEn != fullnameEn || teacher.Mobile != mobile {
		teacher.Fullname = fullname
		teacher.FullnameEn = fullnameEn
		teacher.Mobile = mobile
		if err = w.MongoClient().UpdateTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to update teacher", err, re.ErrorDatabase)
		}
	}
	// 检查时间是否有冲突
	theDay := utils.BeginOfDay(start)
	nextDay := utils.BeginOfTomorrow(start)
	theDayReservations, err := w.MongoClient().GetReservationsBetweenTime(theDay, nextDay)
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
		StartTime:         start,
		EndTime:           end,
		Status:            model.ReservationStatusAvailable,
		InternationalType: teacher.InternationalType,
		TeacherUsername:   teacher.Username,
		TeacherFullname:   teacher.Fullname,
		TeacherFullnameEn: teacher.FullnameEn,
		TeacherMobile:     teacher.Mobile,
		TeacherAddress:    teacher.Address,
		TeacherAddressEn:  teacher.AddressEn,
	}
	if err = w.MongoClient().InsertReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to insert new reservation", err, re.ErrorDatabase)
	}
	return reservation, nil
}

func (w *Workflow) EditReservationByTeacher(reservationId string, startTime string, endTime string, fullname string,
	fullnameEn string, mobile string, userId string, userType int) (*model.Reservation, error) {
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
	} else if fullname == "" && fullnameEn == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	}
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.MongoClient().GetReservationById(reservationId)
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
	if teacher.Fullname != fullname || teacher.FullnameEn != fullnameEn || teacher.Mobile != mobile {
		teacher.Fullname = fullname
		teacher.FullnameEn = fullnameEn
		teacher.Mobile = mobile
		if err = w.MongoClient().UpdateTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to update teacher", err, re.ErrorDatabase)
		}
	}
	// 检查时间是否有冲突
	theDay := utils.BeginOfDay(start)
	nextDay := utils.BeginOfTomorrow(start)
	theDayReservations, err := w.MongoClient().GetReservationsBetweenTime(theDay, nextDay)
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
	reservation.TeacherFullnameEn = teacher.FullnameEn
	reservation.TeacherMobile = teacher.Mobile
	if err = w.MongoClient().UpdateReservation(reservation); err != nil {
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
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return 0, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	removed := 0
	for _, reservationId := range reservationIds {
		reservation, err := w.MongoClient().GetReservationById(reservationId)
		if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted ||
			reservation.TeacherUsername != teacher.Username {
			continue
		}
		reservation.Status = model.ReservationStatusDeleted
		if w.MongoClient().UpdateReservation(reservation) == nil {
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
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return 0, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	canceled := 0
	for _, reservationId := range reservationIds {
		reservation, err := w.MongoClient().GetReservationById(reservationId)
		if err != nil || reservation == nil || reservation.Status != model.ReservationStatusReservated ||
			reservation.StartTime.Before(time.Now()) || reservation.TeacherUsername != teacher.Username {
			continue
		}
		studentFullname, studentMobile := reservation.StudentInfo.Fullname, reservation.StudentInfo.Mobile
		reservation.Status = model.ReservationStatusAvailable
		reservation.StudentInfo = model.StudentInfo{}
		reservation.StudentFeedback = model.StudentFeedback{}
		reservation.TeacherFeedback = model.TeacherFeedback{}
		if w.MongoClient().UpdateReservation(reservation) == nil {
			canceled++
		}
		go w.SendCancelSMS(reservation, studentFullname, studentMobile)
	}
	return canceled, nil
}

func (w *Workflow) MakeReservationByTeacher(reservationId string, fullname string, gender string,
	username string, school string, hometown string, mobile string, email string, experience string,
	problem string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
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
	}

	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
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
	} else if reservation.TeacherUsername != teacher.Username {
		return nil, re.NewRErrorCode("cannot make other teacher's reservation", nil, re.ErrorNotAuthorized)
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

	//send success sms
	go w.SendSuccessSMS(reservation, teacher)
	return reservation, nil
}

func (w *Workflow) GetFeedbackByTeacher(reservationId string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	}
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.MongoClient().GetReservationById(reservationId)
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

func (w *Workflow) SubmitFeedbackByTeacher(reservationId string, teacherFullname string, teacherFullnameEn string,
	teacherUsername string, studentFullname string, problem string, solution string, adviceToCenter string,
	userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	} else if teacherFullname == "" && teacherFullnameEn == "" {
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
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.MongoClient().GetReservationById(reservationId)
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
		TeacherFullname:   teacherFullname,
		TeacherFullnameEn: teacherFullnameEn,
		TeacherUsername:   teacherUsername,
		StudentFullname:   studentFullname,
		Problem:           problem,
		Solution:          solution,
		AdviceToCenter:    adviceToCenter,
	}
	if err = w.MongoClient().UpdateReservation(reservation); err != nil {
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
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	reservation, err := w.MongoClient().GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorViewAvailableReservationStudentInfo)
	} else if reservation.TeacherUsername != teacher.Username {
		return nil, re.NewRErrorCode("cannot get feedback of other one's reservation", nil, re.ErrorTeacherViewOtherReservation)
	}
	return &reservation.StudentInfo, nil
}

func (w *Workflow) UpdateSmsSuffixByTeacher(smsSuffix string, smsSuffixEn string, userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	}
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	teacher.SmsSuffix = smsSuffix
	teacher.SmsSuffixEn = smsSuffixEn
	if err = w.MongoClient().UpdateTeacher(teacher); err != nil {
		return nil, re.NewRErrorCode("fail to update teacher", err, re.ErrorDatabase)
	}
	return teacher, nil
}

func (w *Workflow) SendSMSByTeacher(mobile string, content string, userId string, userType int) error {
	if userId == "" {
		return re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if mobile == "" {
		return re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if content == "" {
		return re.NewRErrorCodeContext("content is empty", nil, re.ErrorMissingParam, "content")
	} else if !utils.IsMobile(mobile) {
		return re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	}

	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}

	return w.sendSMS(mobile, content)
}

func (w *Workflow) SendEmailByTeacher(email string, subject string, body string, userId string, userType int) error {
	if userId == "" {
		return re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	} else if email == "" {
		return re.NewRErrorCodeContext("email is empty", nil, re.ErrorMissingParam, "email")
	} else if subject == "" {
		return re.NewRErrorCodeContext("subject is empty", nil, re.ErrorMissingParam, "subject")
	} else if body == "" {
		return re.NewRErrorCodeContext("body is empty", nil, re.ErrorMissingParam, "body")
	} else if !utils.IsEmail(email) {
		return re.NewRErrorCode("email format is wrong", nil, re.ErrorFormatEmail)
	}

	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}

	return SendEmail(subject, body, []string{email})
}

func (w *Workflow) WrapSimpleTeacher(teacher *model.Teacher) map[string]interface{} {
	var result = make(map[string]interface{})
	if teacher == nil {
		return result
	}
	result["fullname"] = teacher.Fullname
	result["fullname_en"] = teacher.FullnameEn
	result["address"] = teacher.Address
	result["address_en"] = teacher.AddressEn
	result["gender"] = teacher.Gender
	result["gender_en"] = teacher.GenderEn
	result["major"] = teacher.Major
	result["major_en"] = teacher.MajorEn
	result["academic"] = teacher.Academic
	result["academic_en"] = teacher.AcademicEn
	result["aptitude"] = teacher.Aptitude
	result["aptitude_en"] = teacher.AptitudeEn
	result["problem"] = teacher.Problem
	result["problem_en"] = teacher.ProblemEn
	result["international_type"] = teacher.InternationalType
	return result
}

func (w *Workflow) WrapTeacher(teacher *model.Teacher) map[string]interface{} {
	var result = w.WrapSimpleTeacher(teacher)
	if teacher == nil {
		return result
	}
	result["username"] = teacher.Username
	result["mobile"] = teacher.Mobile
	result["sms_suffix"] = teacher.SmsSuffix
	result["sms_suffix_en"] = teacher.SmsSuffixEn
	return result
}
