package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"strings"
	"time"
)

type TeacherLogic struct {
}

// 咨询师添加咨询
func (tl *TeacherLogic) AddReservationByTeacher(startTime string, endTime string, teacherFullname string,
	teacherMobile string, userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if len(startTime) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endTime) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	start, err := time.ParseInLocation(utils.TIME_PATTERN, startTime, utils.Location)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation(utils.TIME_PATTERN, endTime, utils.Location)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertUser(teacher); err != nil {
			return nil, errors.New("数据获取失败")
		}
	}
	// 检查当天的咨询师时间是否有冲突
	theDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, utils.Location)
	nextDay := theDay.AddDate(0, 0, 1)
	theDayReservations, err := models.GetReservationsBetweenTime(theDay, nextDay)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	for _, r := range theDayReservations {
		if strings.EqualFold(r.TeacherUsername, teacher.Username) {
			if (start.After(r.StartTime) && start.Before(r.EndTime)) ||
			(end.After(r.StartTime) && end.Before(r.EndTime)) ||
			(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, errors.New("咨询师时间有冲突")
			}
		}
	}
	reservation, err := models.AddReservation(start, end, teacher.Fullname, teacher.Username, teacher.Mobile)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 咨询师编辑咨询
func (tl *TeacherLogic) EditReservationByTeacher(reservationId string, startTime string, endTime string,
	teacherFullname string, teacherMobile string, userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(startTime) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endTime) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == models.RESERVATED {
		return nil, errors.New("不能编辑已被预约的咨询")
	}
	start, err := time.ParseInLocation(utils.TIME_PATTERN, startTime, utils.Location)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.ParseInLocation(utils.TIME_PATTERN, endTime, utils.Location)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	} else if end.Before(time.Now().In(utils.Location)) {
		return nil, errors.New("不能编辑已过期咨询")
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Username, reservation.TeacherUsername) {
		return nil, errors.New("只能编辑本人开设的咨询")
	}
	if !strings.EqualFold(teacher.Fullname, teacherFullname) || !strings.EqualFold(teacher.Mobile, teacherMobile) {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertUser(teacher); err != nil {
			return nil, errors.New("数据获取失败")
		}
	}
	// 检查当天的咨询师时间是否有冲突
	theDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, utils.Location)
	nextDay := theDay.AddDate(0, 0, 1)
	theDayReservations, err := models.GetReservationsBetweenTime(theDay, nextDay)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	for _, r := range theDayReservations {
		if !strings.EqualFold(r.Id.Hex(), reservation.Id.Hex()) &&
			strings.EqualFold(r.TeacherUsername, teacher.Username) {
			if (start.After(r.StartTime) && start.Before(r.EndTime)) ||
			(end.After(r.StartTime) && end.Before(r.EndTime)) ||
			(!start.After(r.StartTime) && !end.Before(r.EndTime)) {
				return nil, errors.New("咨询师时间有冲突")
			}
		}
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherFullname = teacherFullname
	reservation.TeacherMobile = teacherMobile
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 咨询师删除咨询
func (tl *TeacherLogic) RemoveReservationsByTeacher(reservationIds []string, userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return 0, errors.New("权限不足")
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return 0, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return 0, errors.New("权限不足")
	}
	removed := 0
	for _, reservationId := range reservationIds {
		if reservation, err := models.GetReservationById(reservationId); err == nil && strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
			reservation.Status = models.DELETED
			models.UpsertReservation(reservation)
			removed++
		}
	}
	return removed, nil
}

// 咨询师取消预约
func (tl *TeacherLogic) CancelReservationsByTeacher(reservationIds []string, userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return 0, errors.New("权限不足")
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return 0, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return 0, errors.New("权限不足")
	}
	removed := 0
	for _, reservationId := range reservationIds {
		reseravtion, err := models.GetReservationById(reservationId)
		if err != nil || reseravtion.Status == models.DELETED ||
			!strings.EqualFold(reseravtion.TeacherUsername, teacher.Username) {
			continue
		}
		if reseravtion.Status == models.RESERVATED && reseravtion.StartTime.After(time.Now().In(utils.Location)) {
			reseravtion.Status = models.AVAILABLE
			reseravtion.StudentInfo = models.StudentInfo{}
			reseravtion.StudentFeedback = models.StudentFeedback{}
			reseravtion.TeacherFeedback = models.TeacherFeedback{}
			models.UpsertReservation(reseravtion)
			removed++
		}
	}
	return removed, nil
}

// 咨询师拉取反馈
func (tl *TeacherLogic) GetFeedbackByTeacher(reservationId string, userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().In(utils.Location)) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
		return nil, errors.New("只能反馈本人开设的咨询")
	}
	return reservation, nil
}

// 咨询师提交反馈
func (tl *TeacherLogic) SubmitFeedbackByTeacher(reservationId string, teacherFullname string, teacherId string,
	studentName string, problem string, solution string, adviceToCenter string, userId string,
	userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherId) == 0 {
		return nil, errors.New("咨询师工作证号为空")
	} else if len(studentName) == 0 {
		return nil, errors.New("学生姓名为空")
	} else if len(problem) == 0 {
		return nil, errors.New("咨询问题为空")
	} else if len(solution) == 0 {
		return nil, errors.New("解决方法为空")
	} else if len(adviceToCenter) == 0 {
		return nil, errors.New("工作建议为空")
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().In(utils.Location)) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
		return nil, errors.New("只能反馈本人开设的咨询")
	}
	if reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty() {
		utils.SendFeedbackSMS(reservation)
	}
	reservation.TeacherFeedback = models.TeacherFeedback{
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherId,
		StudentFullname: studentName,
		Problem:         problem,
		Solution:        solution,
		AdviceToCenter:  adviceToCenter,
	}
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 咨询师查看学生信息
func (tl *TeacherLogic) GetStudentInfoByTeacher(reservationId string, userId string, userType models.UserType) (*models.StudentInfo, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,无法查看")
	} else if !strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
		return nil, errors.New("只能查看本人开设的咨询")
	}
	return &reservation.StudentInfo, nil
}
