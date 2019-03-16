package buslogic

import (
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"sort"
	"time"
)

func (w *Workflow) GetReservationsByStudent(language string) ([]*model.Reservation, error) {
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now().AddDate(0, 0, 8)
	reservations, err := w.MongoClient().GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ErrorDatabase)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.ReservationStatusAvailable && r.StartTime.Before(time.Now()) {
			continue
		}
		if (language == "zh_cn" && (r.InternationalType == model.InternationalTypeChinese || r.InternationalType == model.InternationalTypeChinglish)) ||
			(language == "en_us" && r.InternationalType == model.InternationalTypeChinglish) {
			// 过滤支持的语言
			result = append(result, r)
		}
	}
	sort.Sort(model.ByStartTimeOfReservation(result))
	return result, nil
}

func (w *Workflow) GetReservationsByTeacher(userId string, userType int) (*model.Teacher, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeTeacher {
		return nil, nil, re.NewRErrorCode("user is not teacher", nil, re.ErrorNotAuthorized)
	}
	teacher, err := w.MongoClient().GetTeacherById(userId)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.MongoClient().GetReservationsAfterTime(from)
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get reservations", err, re.ErrorDatabase)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.ReservationStatusAvailable && r.StartTime.Before(time.Now()) {
			continue
		} else if r.TeacherUsername == teacher.Username {
			result = append(result, r)
		}
	}
	sort.Sort(model.ByStartTimeOfReservation(result))
	return teacher, result, nil
}

func (w *Workflow) GetReservationsByAdmin(userId string, userType int) (*model.Admin, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	}
	admin, err := w.MongoClient().GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.MongoClient().GetReservationsAfterTime(from)
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get reservations", err, re.ErrorDatabase)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.ReservationStatusAvailable && r.StartTime.Before(time.Now()) {
			continue
		}
		result = append(result, r)
	}
	sort.Sort(model.ByStartTimeOfReservation(result))
	return admin, result, nil
}

func (w *Workflow) GetReservationsMonthlyByAdmin(fromDate string, userId string, userType int) (*model.Admin, []*model.Reservation, error) {
	if userId == "" {
		return nil, nil, re.NewRErrorCode("teacher not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	}
	admin, err := w.MongoClient().GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	from, err := time.ParseInLocation("2006-01-02", fromDate, time.Local)
	if err != nil {
		return nil, nil, re.NewRErrorCodeContext("from date is not valid", err, re.ErrorInvalidParam, "from_date")
	}
	to := from.AddDate(0, 0, 1)
	reservations, err := w.MongoClient().GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, nil, re.NewRErrorCode("fail to get reservations", err, re.ErrorDatabase)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.ReservationStatusAvailable && r.StartTime.Before(time.Now()) {
			continue
		}
		result = append(result, r)
	}
	sort.Sort(model.ByStartTimeOfReservation(result))
	return admin, result, nil
}

// external 将所有咨询移动n天
func (w *Workflow) ShiftReservationTimeInDays(days int) error {
	reservations, err := w.MongoClient().GetAllReservations()
	if err != nil {
		return err
	}
	for _, r := range reservations {
		r.StartTime = r.StartTime.AddDate(0, 0, days)
		r.EndTime = r.EndTime.AddDate(0, 0, days)
		err = w.MongoClient().UpdateReservationWithoutUpdatedTime(r)
		if err != nil {
			log.Errorf("fail to update reservation %+v, err: %+v", r, err)
		}
	}
	return nil
}

type ReservationGroup struct {
	DateStr                   string
	TeacherUsernameMap        map[string]int
	TotalReservationCount     int
	AvailableReservationCount int
	Reservations              []*model.Reservation
}

func (w *Workflow) GroupReservationsForStudent(reservations []*model.Reservation) []*ReservationGroup {
	reservationGroupMap := make(map[string]*ReservationGroup)
	for _, res := range reservations {
		date := res.StartTime.Format("2006-01-02")
		if _, ok := reservationGroupMap[date]; !ok {
			reservationGroupMap[date] = &ReservationGroup{
				DateStr:            date,
				TeacherUsernameMap: make(map[string]int),
				Reservations:       make([]*model.Reservation, 0),
			}
		}
		reservationGroupMap[date].TeacherUsernameMap[res.TeacherUsername]++
		reservationGroupMap[date].TotalReservationCount++
		if res.Status == model.ReservationStatusAvailable {
			reservationGroupMap[date].AvailableReservationCount++
		}
		reservationGroupMap[date].Reservations = append(reservationGroupMap[date].Reservations, res)
	}
	i := 0
	reservationGroups := make([]*ReservationGroup, len(reservationGroupMap))
	for _, resGroup := range reservationGroupMap {
		reservationGroups[i] = resGroup
		i++
	}
	sort.Slice(reservationGroups, func(i, j int) bool {
		return reservationGroups[i].DateStr < reservationGroups[j].DateStr
	})
	return reservationGroups
}

func (w *Workflow) WrapReservationGroupForStudent(resGroup *ReservationGroup) map[string]interface{} {
	var result = make(map[string]interface{})
	if resGroup == nil {
		return result
	}
	result["date"] = resGroup.DateStr
	result["teacher_num"] = len(resGroup.TeacherUsernameMap)
	result["total_reservation_count"] = resGroup.TotalReservationCount
	result["available_reservation_count"] = resGroup.AvailableReservationCount
	var array = make([]interface{}, 0)
	for _, res := range resGroup.Reservations {
		array = append(array, w.WrapSimpleReservation(res))
	}
	result["reservations"] = array
	return result
}

func (w *Workflow) WrapSimpleReservation(reservation *model.Reservation) map[string]interface{} {
	var result = make(map[string]interface{})
	if reservation == nil {
		return result
	}
	result["id"] = reservation.Id.Hex()
	result["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	result["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	result["status"] = reservation.Status
	result["international_type"] = reservation.InternationalType
	if reservation.Status == model.ReservationStatusReservated && reservation.StartTime.Before(time.Now()) {
		result["status"] = model.ReservationStatusFeedback
	}
	result["teacher_fullname"] = reservation.TeacherFullname
	result["teacher_fullname_en"] = reservation.TeacherFullnameEn
	result["teacher_address"] = reservation.TeacherAddress
	result["teacher_address_en"] = reservation.TeacherAddressEn
	return result
}

func (w *Workflow) WrapReservation(reservation *model.Reservation) map[string]interface{} {
	var result = w.WrapSimpleReservation(reservation)
	if reservation == nil {
		return result
	}
	result["teacher_username"] = reservation.TeacherUsername
	result["teacher_mobile"] = reservation.TeacherMobile
	return result
}

func (w *Workflow) WrapReservationStudentFeedback(studentFeedback model.StudentFeedback) map[string]interface{} {
	var result = make(map[string]interface{})
	result["fullname"] = studentFeedback.Fullname
	result["problem"] = studentFeedback.Problem
	result["choices"] = studentFeedback.Choices
	result["score"] = studentFeedback.Score
	result["feedback"] = studentFeedback.Feedback
	return result
}

func (w *Workflow) WrapReservationTeacherFeedback(teacherFeedback model.TeacherFeedback) map[string]interface{} {
	var result = make(map[string]interface{})
	result["teacher_fullname"] = teacherFeedback.TeacherFullname
	result["teacher_fullname_en"] = teacherFeedback.TeacherFullnameEn
	result["teacher_username"] = teacherFeedback.TeacherUsername
	result["student_fullname"] = teacherFeedback.StudentFullname
	result["problem"] = teacherFeedback.Problem
	result["solution"] = teacherFeedback.Solution
	result["advice_to_center"] = teacherFeedback.AdviceToCenter
	return result
}
