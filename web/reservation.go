package web

import (
	"github.com/mijia/sweb/form"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	"github.com/shudiwsh2009/reservation_thxx_go/service"
	"net/http"
)

type ReservationController struct {
	BaseMuxController
}

const (
	kStudentApiBaseUrl = "/api/student"
	kTeacherApiBaseUrl = "/api/teacher"
	kAdminApiBaseUrl   = "/api/admin"
)

func (rc *ReservationController) MuxHandlers(m JsonMuxer) {
	m.GetJson(kStudentApiBaseUrl+"/reservation/view", "ViewReservationsByStudent", FakeInjection(rc.ViewReservationsByStudent))
	m.GetJson(kStudentApiBaseUrl+"/reservation/view/group", "ViewGroupedReservationsByStudent", FakeInjection(rc.ViewGroupedReservationsByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/make", "MakeReservationByStudent", FakeInjection(rc.MakeReservationByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByStudent", FakeInjection(rc.GetFeedbackByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByStudent", FakeInjection(rc.SubmitFeedbackByStudent))
	m.PostJson(kStudentApiBaseUrl+"/reservation/teacher/get", "GetReservationTeacherInfoByStudent", FakeInjection(rc.GetReservationTeacherInfoByStudent))

	m.GetJson(kTeacherApiBaseUrl+"/reservation/view", "ViewReservationsByTeacher", RoleCookieInjection(rc.ViewReservationsByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/add", "AddReservationByTeacher", RoleCookieInjection(rc.AddReservationByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/edit", "EditReservationByTeacher", RoleCookieInjection(rc.EditReservationByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/remove", "RemoveReservationsByTeacher", RoleCookieInjection(rc.RemoveReservationsByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/cancel", "CancelReservationsByTeacher", RoleCookieInjection(rc.CancelReservationsByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/make", "MakeReservationByTeacher", RoleCookieInjection(rc.MakeReservationByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByTeacher", RoleCookieInjection(rc.GetFeedbackByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByTeacher", RoleCookieInjection(rc.SubmitFeedbackByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/reservation/student/get", "GetReservationStudentInfoByTeacher", RoleCookieInjection(rc.GetReservationStudentInfoByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/sms_suffix/update", "UpdateSmsSuffixByTeacher", RoleCookieInjection(rc.UpdateSmsSuffixByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/sms/send", "SendSMSByTeacher", RoleCookieInjection(rc.SendSMSByTeacher))
	m.PostJson(kTeacherApiBaseUrl+"/email/send", "SendEmailByTeacher", RoleCookieInjection(rc.SendEmailByTeacher))

	m.GetJson(kAdminApiBaseUrl+"/reservation/view", "ViewReservationsByAdmin", RoleCookieInjection(rc.ViewReservationsByAdmin))
	m.GetJson(kAdminApiBaseUrl+"/reservation/view/monthly", "ViewReservationsMonthlyByAdmin", RoleCookieInjection(rc.ViewReservationsMonthlyByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/add", "AddReservationByAdmin", RoleCookieInjection(rc.AddReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/edit", "EditReservationByAdmin", RoleCookieInjection(rc.EditReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/remove", "RemoveReservationsByAdmin", RoleCookieInjection(rc.RemoveReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/cancel", "CancelReservationsByAdmin", RoleCookieInjection(rc.CancelReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/make", "MakeReservationByAdmin", RoleCookieInjection(rc.MakeReservationByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/feedback/get", "GetFeedbackByAdmin", RoleCookieInjection(rc.GetFeedbackByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/feedback/submit", "SubmitFeedbackByAdmin", RoleCookieInjection(rc.SubmitFeedbackByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/student/get", "GetReservationStudentInfoByAdmin", RoleCookieInjection(rc.GetReservationStudentInfoByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/search", "SearchTeacherByAdmin", RoleCookieInjection(rc.SearchTeacherByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/get", "GetTeacherInfoByAdmin", RoleCookieInjection(rc.GetTeacherInfoByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/teacher/edit", "EditTeacherInfoByAdmin", RoleCookieInjection(rc.EditTeacherInfoByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/export", "ExportReservationsByAdmin", RoleCookieInjection(rc.ExportReservationsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/reservation/export/arrangements", "ExportReservationArrangementsByAdmin", RoleCookieInjection(rc.ExportReservationArrangementsByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/sms/send", "SendSMSByAdmin", RoleCookieInjection(rc.SendSMSByAdmin))
	m.PostJson(kAdminApiBaseUrl+"/email/send", "SendEmailByAdmin", RoleCookieInjection(rc.SendEmailByAdmin))
}

func (rc *ReservationController) ViewReservationsByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	language := form.ParamString(r, "language", "zh_cn")

	var result = make(map[string]interface{})

	reservations, err := service.Workflow().GetReservationsByStudent(language)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	var array = make([]interface{}, len(reservations))
	for i, res := range reservations {
		array[i] = service.Workflow().WrapSimpleReservation(res)
	}
	result["reservations"] = array

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ViewGroupedReservationsByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	language := form.ParamString(r, "language", "zh_cn")

	var result = make(map[string]interface{})

	reservations, err := service.Workflow().GetReservationsByStudent(language)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	var array = make([]interface{}, len(reservations))
	for i, res := range reservations {
		array[i] = service.Workflow().WrapSimpleReservation(res)
	}
	result["reservations"] = array

	reservationGroups := service.Workflow().GroupReservationsForStudent(reservations)
	var groupArray = make([]interface{}, len(reservationGroups))
	for i, resGroup := range reservationGroups {
		groupArray[i] = service.Workflow().WrapReservationGroupForStudent(resGroup)
	}
	result["reservation_groups"] = groupArray

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) MakeReservationByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	fullname := form.ParamString(r, "fullname", "")
	gender := form.ParamString(r, "gender", "")
	username := form.ParamString(r, "username", "")
	school := form.ParamString(r, "school", "")
	hometown := form.ParamString(r, "hometown", "")
	mobile := form.ParamString(r, "mobile", "")
	email := form.ParamString(r, "email", "")
	experience := form.ParamString(r, "experience", "")
	problem := form.ParamString(r, "problem", "")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().MakeReservationByStudent(reservationId, fullname, gender, username,
		school, hometown, mobile, email, experience, problem)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapSimpleReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	username := form.ParamString(r, "username", "")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().GetFeedbackByStudent(reservationId, username)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	feedbackJson := service.Workflow().WrapReservationStudentFeedback(reservation.StudentFeedback)
	if reservation.StudentFeedback.Fullname == "" {
		feedbackJson["fullname"] = reservation.StudentInfo.Fullname
	}
	if reservation.StudentFeedback.Problem == "" {
		feedbackJson["problem"] = reservation.StudentInfo.Problem
	}
	result["feedback"] = feedbackJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	fullname := form.ParamString(r, "fullname", "")
	problem := form.ParamString(r, "problem", "")
	choices := form.ParamString(r, "choices", "")
	score := form.ParamString(r, "score", "")
	feedback := form.ParamString(r, "feedback", "")
	username := form.ParamString(r, "username", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().SubmitFeedbackByStudent(reservationId, fullname, problem, choices, score, feedback, username)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetReservationTeacherInfoByStudent(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().GetReservationTeacherInfoByStudent(reservationId)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapSimpleTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ViewReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	teacher, reservations, err := service.Workflow().GetReservationsByTeacher(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
		array = append(array, resJson)
	}
	result["reservations"] = array

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) AddReservationByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	startTime := form.ParamString(r, "start_time", "")
	endTime := form.ParamString(r, "end_time", "")
	fullname := form.ParamString(r, "fullname", "")
	fullnameEn := form.ParamString(r, "fullname_en", "")
	mobile := form.ParamString(r, "mobile", "")
	location := form.ParamInt(r, "location", 0)

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().AddReservationByTeacher(startTime, endTime, fullname, fullnameEn, mobile, location, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) EditReservationByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	startTime := form.ParamString(r, "start_time", "")
	endTime := form.ParamString(r, "end_time", "")
	fullname := form.ParamString(r, "fullname", "")
	fullnameEn := form.ParamString(r, "fullname_en", "")
	mobile := form.ParamString(r, "mobile", "")
	location := form.ParamInt(r, "location", 0)

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().EditReservationByTeacher(reservationId, startTime, endTime, fullname, fullnameEn, mobile, location, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) RemoveReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationIds := []string(r.Form["reservation_ids"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().RemoveReservationsByTeacher(reservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["removed_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) CancelReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationIds := []string(r.Form["reservation_ids"])

	var result = make(map[string]interface{})

	canceled, err := service.Workflow().CancelReservationsByTeacher(reservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["canceled_count"] = canceled

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) MakeReservationByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	fullname := form.ParamString(r, "fullname", "")
	gender := form.ParamString(r, "gender", "")
	username := form.ParamString(r, "username", "")
	school := form.ParamString(r, "school", "")
	hometown := form.ParamString(r, "hometown", "")
	mobile := form.ParamString(r, "mobile", "")
	email := form.ParamString(r, "email", "")
	experience := form.ParamString(r, "experience", "")
	problem := form.ParamString(r, "problem", "")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().MakeReservationByTeacher(reservationId, fullname, gender, username,
		school, hometown, mobile, email, experience, problem, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().GetFeedbackByTeacher(reservationId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	feedbackJson := service.Workflow().WrapReservationTeacherFeedback(reservation.TeacherFeedback)
	if reservation.TeacherFeedback.TeacherFullname == "" {
		feedbackJson["teacher_fullname"] = reservation.TeacherFullname
	}
	if reservation.TeacherFeedback.TeacherFullnameEn == "" {
		feedbackJson["teacher_fullname_en"] = reservation.TeacherFullnameEn
	}
	if reservation.TeacherFeedback.TeacherUsername == "" {
		feedbackJson["teacher_username"] = reservation.TeacherUsername
	}
	if reservation.TeacherFeedback.StudentFullname == "" {
		feedbackJson["student_fullname"] = reservation.StudentInfo.Fullname
	}
	result["feedback"] = feedbackJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherFullnameEn := form.ParamString(r, "teacher_fullname_en", "")
	studentFullname := form.ParamString(r, "student_fullname", "")
	problem := form.ParamString(r, "problem", "")
	solution := form.ParamString(r, "solution", "")
	adviceToCenter := form.ParamString(r, "advice_to_center", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().SubmitFeedbackByTeacher(reservationId, teacherFullname, teacherFullnameEn,
		teacherUsername, studentFullname, problem, solution, adviceToCenter, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetReservationStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")

	var result = make(map[string]interface{})

	studentInfo, err := service.Workflow().GetReservationStudentInfoByTeacher(reservationId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["student"] = service.Workflow().WrapStudenInfo(studentInfo)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) UpdateSmsSuffixByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	smsSuffix := form.ParamString(r, "sms_suffix", "")
	smsSuffixEn := form.ParamString(r, "sms_suffix_en", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().UpdateSmsSuffixByTeacher(smsSuffix, smsSuffixEn, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SendSMSByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	mobile := form.ParamString(r, "mobile", "")
	content := form.ParamString(r, "content", "")

	var result = make(map[string]interface{})

	err := service.Workflow().SendSMSByTeacher(mobile, content, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SendEmailByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	email := form.ParamString(r, "email", "")
	subject := form.ParamString(r, "subject", "")
	body := form.ParamString(r, "body", "")

	var result = make(map[string]interface{})

	err := service.Workflow().SendEmailByTeacher(email, subject, body, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ViewReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	_, reservations, err := service.Workflow().GetReservationsByAdmin(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
		array = append(array, resJson)
	}
	result["reservations"] = array

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ViewReservationsMonthlyByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := form.ParamString(r, "from_date", "")

	var result = make(map[string]interface{})

	_, reservations, err := service.Workflow().GetReservationsMonthlyByAdmin(fromDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := service.Workflow().WrapReservation(res)
		array = append(array, resJson)
	}
	result["reservations"] = array

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) AddReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	startTime := form.ParamString(r, "start_time", "")
	endTime := form.ParamString(r, "end_time", "")
	username := form.ParamString(r, "username", "")
	fullname := form.ParamString(r, "fullname", "")
	fullnameEn := form.ParamString(r, "fullname_en", "")
	mobile := form.ParamString(r, "mobile", "")
	address := form.ParamString(r, "address", "")
	addressEn := form.ParamString(r, "address_en", "")
	internationalType := form.ParamInt(r, "international_type", 0)
	location := form.ParamInt(r, "location", 0)

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().AddReservationByAdmin(startTime, endTime, username, fullname, fullnameEn, mobile, address, addressEn, internationalType, location, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) EditReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	startTime := form.ParamString(r, "start_time", "")
	endTime := form.ParamString(r, "end_time", "")
	username := form.ParamString(r, "username", "")
	fullname := form.ParamString(r, "fullname", "")
	fullnameEn := form.ParamString(r, "fullname_en", "")
	mobile := form.ParamString(r, "mobile", "")
	address := form.ParamString(r, "address", "")
	addressEn := form.ParamString(r, "address_en", "")
	internationalType := form.ParamInt(r, "international_type", 0)
	location := form.ParamInt(r, "location", 0)

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().EditReservationByAdmin(reservationId, startTime, endTime, username,
		fullname, fullnameEn, mobile, address, addressEn, internationalType, location, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) RemoveReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationIds := []string(r.Form["reservation_ids"])

	var result = make(map[string]interface{})

	removed, err := service.Workflow().RemoveReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["removed_count"] = removed

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) CancelReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationIds := []string(r.Form["reservation_ids"])

	var result = make(map[string]interface{})

	canceled, err := service.Workflow().CancelReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["canceled_count"] = canceled

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) MakeReservationByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	fullname := form.ParamString(r, "fullname", "")
	gender := form.ParamString(r, "gender", "")
	username := form.ParamString(r, "username", "")
	school := form.ParamString(r, "school", "")
	hometown := form.ParamString(r, "hometown", "")
	mobile := form.ParamString(r, "mobile", "")
	email := form.ParamString(r, "email", "")
	experience := form.ParamString(r, "experience", "")
	problem := form.ParamString(r, "problem", "")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().MakeReservationByAdmin(reservationId, fullname, gender, username,
		school, hometown, mobile, email, experience, problem, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["reservation"] = service.Workflow().WrapReservation(reservation)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")

	var result = make(map[string]interface{})

	reservation, err := service.Workflow().GetFeedbackByAdmin(reservationId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	feedbackJson := service.Workflow().WrapReservationTeacherFeedback(reservation.TeacherFeedback)
	if reservation.TeacherFeedback.TeacherFullname == "" {
		feedbackJson["teacher_fullname"] = reservation.TeacherFullname
	}
	if reservation.TeacherFeedback.TeacherFullnameEn == "" {
		feedbackJson["teacher_fullname_en"] = reservation.TeacherFullnameEn
	}
	if reservation.TeacherFeedback.TeacherUsername == "" {
		feedbackJson["teacher_username"] = reservation.TeacherUsername
	}
	if reservation.TeacherFeedback.StudentFullname == "" {
		feedbackJson["student_fullname"] = reservation.StudentInfo.Fullname
	}
	result["feedback"] = feedbackJson

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SubmitFeedbackByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")
	teacherUsername := form.ParamString(r, "teacher_username", "")
	teacherFullname := form.ParamString(r, "teacher_fullname", "")
	teacherFullnameEn := form.ParamString(r, "teacher_fullname_en", "")
	studentFullname := form.ParamString(r, "student_fullname", "")
	problem := form.ParamString(r, "problem", "")
	solution := form.ParamString(r, "solution", "")
	adviceToCenter := form.ParamString(r, "advice_to_center", "")

	var result = make(map[string]interface{})

	_, err := service.Workflow().SubmitFeedbackByAdmin(reservationId, teacherFullname, teacherFullnameEn,
		teacherUsername, studentFullname, problem, solution, adviceToCenter, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetReservationStudentInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationId := form.ParamString(r, "reservation_id", "")

	var result = make(map[string]interface{})

	studentInfo, err := service.Workflow().GetReservationStudentInfoByAdmin(reservationId, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["student"] = service.Workflow().WrapStudenInfo(studentInfo)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SearchTeacherByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	fullname := form.ParamString(r, "fullname", "")
	fullnameEn := form.ParamString(r, "fullname_en", "")
	mobile := form.ParamString(r, "mobile", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().SearchTeacherByAdmin(fullname, fullnameEn, username, mobile, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) GetTeacherInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	username := form.ParamString(r, "username", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().GetTeacherInfoByAdmin(username, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["teacher"] = service.Workflow().WrapTeacher(teacher)

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) EditTeacherInfoByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	fullname := form.ParamString(r, "fullname", "")
	fullnameEn := form.ParamString(r, "fullname_en", "")
	gender := form.ParamString(r, "gender", "")
	genderEn := form.ParamString(r, "gender_en", "")
	major := form.ParamString(r, "major", "")
	majorEn := form.ParamString(r, "major_en", "")
	academic := form.ParamString(r, "academic", "")
	academicEn := form.ParamString(r, "academic_en", "")
	aptitude := form.ParamString(r, "aptitude", "")
	aptitudeEn := form.ParamString(r, "aptitude_en", "")
	problem := form.ParamString(r, "problem", "")
	problemEn := form.ParamString(r, "problem_en", "")
	internationalType := form.ParamInt(r, "international_type", model.InternationalTypeChinese)

	var result = make(map[string]interface{})

	_, err := service.Workflow().EditTeacherInfoByAdmin(username, fullname, fullnameEn, gender, genderEn, major, majorEn,
		academic, academicEn, aptitude, aptitudeEn, problem, problemEn, internationalType, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportReservationsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	reservationIds := []string(r.Form["reservation_ids"])

	var result = make(map[string]interface{})

	path, err := service.Workflow().ExportReservationsByAdmin(reservationIds, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["url"] = "/" + path

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) ExportReservationArrangementsByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	fromDate := form.ParamString(r, "from_date", "")

	var result = make(map[string]interface{})

	url, err := service.Workflow().ExportReservationArrangementsByAdmin(fromDate, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["url"] = "/" + url

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SendSMSByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	mobile := form.ParamString(r, "mobile", "")
	content := form.ParamString(r, "content", "")

	var result = make(map[string]interface{})

	err := service.Workflow().SendSMSByAdmin(mobile, content, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (rc *ReservationController) SendEmailByAdmin(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	email := form.ParamString(r, "email", "")
	subject := form.ParamString(r, "subject", "")
	body := form.ParamString(r, "body", "")

	var result = make(map[string]interface{})

	err := service.Workflow().SendEmailByAdmin(email, subject, body, userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}

	return http.StatusOK, wrapJsonOk(result)
}
