package buslogic

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"sort"
	"time"
)

func (w *Workflow) AddReservationByAdmin(startTime string, endTime string, username string,
	fullname string, mobile string, address string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	} else if startTime == "" {
		return nil, re.NewRErrorCodeContext("start_time is empty", nil, re.ErrorMissingParam, "start_time")
	} else if endTime == "" {
		return nil, re.NewRErrorCodeContext("end_time is empty", nil, re.ErrorMissingParam, "end_time")
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if address == "" {
		return nil, re.NewRErrorCodeContext("address is empty", nil, re.ErrorMissingParam, "address")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
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
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	} else if teacher == nil || teacher.UserType != model.UserTypeTeacher {
		teacher = &model.Teacher{
			Username: username,
			Password: TeacherDefaultPassword,
			UserType: model.UserTypeTeacher,
			Fullname: fullname,
			Mobile:   mobile,
			Address:  address,
		}
		if err = w.mongoClient.InsertTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to insert teacher", err, re.ErrorDatabase)
		}
	} else if teacher.Fullname != fullname || teacher.Mobile != mobile || teacher.Address != address {
		teacher.Fullname = fullname
		teacher.Mobile = mobile
		teacher.Address = address
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

func (w *Workflow) EditReservationByAdmin(reservationId string, startTime string, endTime string, username string,
	fullname string, mobile string, address string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	} else if startTime == "" {
		return nil, re.NewRErrorCodeContext("start_time is empty", nil, re.ErrorMissingParam, "start_time")
	} else if endTime == "" {
		return nil, re.NewRErrorCodeContext("end_time is empty", nil, re.ErrorMissingParam, "end_time")
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	} else if mobile == "" {
		return nil, re.NewRErrorCodeContext("mobile is empty", nil, re.ErrorMissingParam, "mobile")
	} else if address == "" {
		return nil, re.NewRErrorCodeContext("address is empty", nil, re.ErrorMissingParam, "address")
	} else if !utils.IsMobile(mobile) {
		return nil, re.NewRErrorCode("mobile format is wrong", nil, re.ErrorFormatMobile)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("cannot get reservation", err, re.ErrorDatabase)
	} else if reservation.Status == model.ReservationStatusReservated {
		return nil, re.NewRErrorCode("cannot edit reservated reservation", nil, re.ErrorEditReservatedReservation)
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
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	} else if teacher == nil || teacher.UserType != model.UserTypeTeacher {
		teacher = &model.Teacher{
			Username: username,
			Password: TeacherDefaultPassword,
			UserType: model.UserTypeTeacher,
			Fullname: fullname,
			Mobile:   mobile,
			Address:  address,
		}
		if err = w.mongoClient.InsertTeacher(teacher); err != nil {
			return nil, re.NewRErrorCode("fail to insert teacher", err, re.ErrorDatabase)
		}
	} else if teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("teacher has wrong user type", nil, re.ErrorDatabase)
	} else if teacher.Fullname != fullname || teacher.Mobile != mobile || teacher.Address != address {
		teacher.Fullname = fullname
		teacher.Mobile = mobile
		teacher.Address = address
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
	// 更新咨询
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherUsername = teacher.Username
	reservation.TeacherFullname = teacher.Fullname
	reservation.TeacherMobile = teacher.Mobile
	reservation.TeacherAddress = teacher.Address
	if err = w.mongoClient.UpdateReservation(reservation); err != nil {
		return nil, re.NewRErrorCode("fail to update reservation", err, re.ErrorDatabase)
	}
	return reservation, nil
}

func (w *Workflow) RemoveReservationsByAdmin(reservationIds []string, userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return 0, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return 0, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	removed := 0
	for _, reservationId := range reservationIds {
		reservation, err := w.mongoClient.GetReservationById(reservationId)
		if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
			continue
		}
		reservation.Status = model.ReservationStatusDeleted
		if w.mongoClient.UpdateReservation(reservation) == nil {
			removed++
		}
	}
	return removed, nil
}

func (w *Workflow) CancelReservationsByAdmin(reservationIds []string, userId string, userType int) (int, error) {
	if userId == "" {
		return 0, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return 0, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return 0, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	canceled := 0
	for _, reservationId := range reservationIds {
		reservation, err := w.mongoClient.GetReservationById(reservationId)
		if err != nil || reservation == nil || reservation.Status != model.ReservationStatusReservated ||
			reservation.StartTime.Before(time.Now()) {
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

func (w *Workflow) GetFeedbackByAdmin(reservationId string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ErrorFeedbackFutureReservation)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorFeedbackAvailableReservation)
	}
	return reservation, nil
}

func (w *Workflow) SubmitFeedbackByAdmin(reservationId string, teacherFullname string, teacherUsername string,
	studentFullname string, problem string, solution string, adviceToCenter string, userId string, userType int) (*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
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
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.StartTime.After(time.Now()) {
		return nil, re.NewRErrorCode("cannot get feedback of future reservation", nil, re.ErrorFeedbackFutureReservation)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorFeedbackAvailableReservation)
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

func (w *Workflow) GetReservationStudentInfoByAdmin(reservationId string, userId string, userType int) (*model.StudentInfo, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	} else if reservationId == "" {
		return nil, re.NewRErrorCodeContext("reservation_id is empty", nil, re.ErrorMissingParam, "reservation_id")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	reservation, err := w.mongoClient.GetReservationById(reservationId)
	if err != nil || reservation == nil || reservation.Status == model.ReservationStatusDeleted {
		return nil, re.NewRErrorCode("fail to get reservation", err, re.ErrorDatabase)
	} else if reservation.Status == model.ReservationStatusAvailable {
		return nil, re.NewRErrorCode("cannot get feedback of available reservation", nil, re.ErrorViewAvailableReservationStudentInfo)
	}
	return &reservation.StudentInfo, nil
}

func (w *Workflow) SearchTeacherByAdmin(fullname string, username string, mobile string,
	userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	if fullname != "" {
		teacher, err := w.mongoClient.GetTeacherByFullname(fullname)
		if err == nil && teacher != nil && teacher.UserType == model.UserTypeTeacher {
			return teacher, nil
		}
	}
	if username != "" {
		teacher, err := w.mongoClient.GetTeacherByUsername(username)
		if err == nil && teacher != nil && teacher.UserType == model.UserTypeTeacher {
			return teacher, nil
		}
	}
	if mobile != "" {
		teacher, err := w.mongoClient.GetTeacherByMobile(mobile)
		if err == nil && teacher != nil && teacher.UserType == model.UserTypeTeacher {
			return teacher, nil
		}
	}
	return nil, re.NewRErrorCode("fail to search teacher", nil, re.ErrorNoUser)
}

func (w *Workflow) GetTeacherInfoByAdmin(username string, userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorNoUser)
	}
	return teacher, nil
}

func (w *Workflow) EditTeacherInfoByAdmin(username string, fullname string, gender string, major string,
	academic string, aptitude string, problem string, userId string, userType int) (*model.Teacher, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	} else if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if fullname == "" {
		return nil, re.NewRErrorCodeContext("fullname is empty", nil, re.ErrorMissingParam, "fullname")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(username)
	if err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
	}
	teacher.Fullname = fullname
	teacher.Gender = gender
	teacher.Major = major
	teacher.Academic = academic
	teacher.Aptitude = aptitude
	teacher.Problem = problem
	if err = w.mongoClient.UpdateTeacher(teacher); err != nil {
		return nil, re.NewRErrorCode("fail to update teacher", err, re.ErrorDatabase)
	}
	return teacher, nil
}

func (w *Workflow) ExportReservationsByAdmin(reservationIds []string, userId string, userType int) (string, error) {
	if userId == "" {
		return "", re.NewRErrorCode("admin not login", nil, re.ErrorNoLogin)
	} else if userType != model.UserTypeAdmin {
		return "", re.NewRErrorCode("user is not admin", nil, re.ErrorNotAuthorized)
	}
	reservations := make([]*model.Reservation, 0)
	for _, rId := range reservationIds {
		reservation, err := w.mongoClient.GetReservationById(rId)
		if err == nil && reservation != nil && reservation.Status != model.ReservationStatusDeleted {
			reservations = append(reservations, reservation)
		}
	}
	sort.Sort(ByStartTimeOfReservation(reservations))
	if len(reservations) == 0 {
		return "", re.NewRErrorCode("no exportable reservations", nil, re.ErrorAdminNoExportableReservations)
	}
	path := filepath.Join(utils.ExportFolder, fmt.Sprintf("export_%s.xlsx", time.Now().Format("20060102")))
	if err := w.ExportReservationsToFile(reservations, path); err != nil {
		return "", re.NewRErrorCode("fail to export reservations", err, re.ErrorAdminExportReservationFailure)
	}
	return path, nil
}

func (w *Workflow) ExportReservationsToFile(reservations []*model.Reservation, path string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	file, err = xlsx.OpenFile(filepath.Join(utils.ExportFolder, "export_template.xlsx"))
	if err != nil {
		return re.NewRError("fail to oepn export template", err)
	}
	sheet = file.Sheet["export"]
	if sheet == nil {
		return re.NewRError("fail to open sheet", err)
	}
	for _, res := range reservations {
		row = sheet.AddRow()
		// 学生申请表
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Fullname)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Gender)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Username)
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.School)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Hometown)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Mobile)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Email)
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Problem)
		// 预约信息
		cell = row.AddCell()
		cell.SetString(res.TeacherFullname)
		cell = row.AddCell()
		cell.SetString(res.StartTime.Format("2006-01-02"))
		// 咨询师反馈表
		cell = row.AddCell()
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.TeacherFeedback.Problem)
		cell = row.AddCell()
		cell.SetString(res.TeacherFeedback.Solution)
		// 学生反馈表
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.StudentFeedback.Score)
		cell = row.AddCell()
		cell.SetString(res.StudentFeedback.Feedback)
		if res.StudentFeedback.Choices != "" {
			for i := 0; i < len(res.StudentFeedback.Choices); i++ {
				cell = row.AddCell()
				switch res.StudentFeedback.Choices[i] {
				case 'A':
					cell.SetString("非常同意")
				case 'B':
					cell.SetString("一般")
				case 'C':
					cell.SetString("不同意")
				default:
				}
			}
		}
	}
	err = file.Save(path)
	if err != nil {
		return re.NewRError("fail to save to file", err)
	}
	return nil
}
