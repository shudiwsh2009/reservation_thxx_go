package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/data"
	"github.com/shudiwsh2009/reservation_thxx_go/domain"
	"github.com/shudiwsh2009/reservation_thxx_go/util"
	"strings"
	"time"
	"github.com/shudiwsh2009/reservation_thxx_go/sms"
)

type TeacherLogic struct {
}

// 咨询师添加咨询
func (tl *TeacherLogic) AddReservationByTeacher(startTime string, endTime string, teacherFullname string,
	teacherMobile string, username string, userType domain.UserType) (*domain.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != domain.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(startTime, "") {
		return nil, errors.New("开始时间为空")
	} else if strings.EqualFold(endTime, "") {
		return nil, errors.New("结束时间为空")
	} else if strings.EqualFold(teacherFullname, "") {
		return nil, errors.New("咨询师姓名为空")
	} else if strings.EqualFold(teacherMobile, "") {
		return nil, errors.New("咨询师手机号为空")
	} else if !util.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	start, err := time.Parse(util.TIME_PATTERN, startTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.Parse(util.TIME_PATTERN, endTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}
	teacher, err := data.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != domain.TEACHER {
		return nil, errors.New("权限不足")
	}
	teacher.Fullname = teacherFullname
	teacher.Mobile = teacherMobile
	if err = data.UpsertUser(teacher); err != nil {
		return nil, errors.New("数据获取失败")
	}
	reservation, err := data.AddReservation(start, end, teacher.Fullname, teacher.Username, teacher.Mobile)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation
}

// 咨询师编辑咨询
func (tl *TeacherLogic) EditReservationByTeacher(reservationId string, startTime string, endTime string,
	teacherFullname string, teacherMobile string, username string, userType domain.UserType) (*domain.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != domain.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(startTime, "") {
		return nil, errors.New("开始时间为空")
	} else if strings.EqualFold(endTime, "") {
		return nil, errors.New("结束时间为空")
	} else if strings.EqualFold(teacherFullname, "") {
		return nil, errors.New("咨询师姓名为空")
	} else if strings.EqualFold(teacherMobile, "") {
		return nil, errors.New("咨询师手机号为空")
	} else if !util.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	reservation, err := data.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == domain.RESERVATED {
		return nil, errors.New("不能编辑已被预约的咨询")
	}
	start, err := time.Parse(util.TIME_PATTERN, startTime)
	if err != nil {
		return nil, errors.New("开始时间格式错误")
	}
	end, err := time.Parse(util.TIME_PATTERN, endTime)
	if err != nil {
		return nil, errors.New("结束时间格式错误")
	}
	if start.After(end) {
		return nil, errors.New("开始时间不能晚于结束时间")
	}
	teacher, err := data.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != domain.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Username, reservation.TeacherUsername) {
		return nil, errors.New("只能编辑本人开设的咨询")
	}
	teacher.Fullname = teacherFullname
	teacher.Mobile = teacherMobile
	if err = data.UpsertUser(teacher); err != nil {
		return nil, errors.New("数据获取失败")
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherFullname = teacherFullname
	reservation.TeacherMobile = teacherMobile
	if err = data.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation
}

// 咨询师删除咨询
func (tl *TeacherLogic) RemoveReservationsByTeacher(reservationIds []string, username string, userType domain.UserType) error {
	if strings.EqualFold(username, "") {
		return errors.New("请先登录")
	} else if userType != domain.TEACHER {
		return errors.New("权限不足")
	} else if reservationIds == nil {
		return errors.New("咨询Id列表为空")
	}
	teacher, err := data.GetUserByUsername(username)
	if err != nil {
		return errors.New("咨询师账户失效")
	} else if teacher.UserType != domain.TEACHER {
		return errors.New("权限不足")
	}
	for _, reservationId := range reservationIds {
		if reservation, err := data.GetReservationById(reservationId); err == nil && strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
			reservation.Status = domain.DELETED
			data.UpsertReservation(reservation)
		}
	}
	return nil
}

// 咨询师取消预约
func (tl *TeacherLogic) CancelReservationsByTeacher(reservationIds []string, username string, userType domain.UserType) error {
	if strings.EqualFold(username, "") {
		return errors.New("请先登录")
	} else if userType != domain.TEACHER {
		return errors.New("权限不足")
	} else if reservationIds == nil {
		return errors.New("咨询Id列表为空")
	}
	teacher, err := data.GetUserByUsername(username)
	if err != nil {
		return errors.New("咨询师账户失效")
	} else if teacher.UserType != domain.TEACHER {
		return errors.New("权限不足")
	}
	for _, reservationId := range reservationIds {
		reseravtion, err := data.GetReservationById(reservationId)
		if err != nil || reseravtion.Status == domain.DELETED ||
			!strings.EqualFold(reseravtion.TeacherUsername, teacher.Username) {
			continue
		}
		if reseravtion.Status == domain.RESERVATED && reseravtion.StartTime.After(time.Now().Local()) {
			reseravtion.Status = domain.AVAILABLE
			reseravtion.StudentInfo = domain.StudentInfo{}
			reseravtion.StudentFeedback = domain.StudentFeedback{}
			reseravtion.TeacherFeedback = domain.TeacherFeedback{}
			data.UpsertReservation(reseravtion)
		}
	}
	return nil
}

// 咨询师拉取反馈
func (tl *TeacherLogic) GetFeedbackByTeacher(reservationId string, username string, userType domain.UserType) (*domain.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != domain.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	}
	teacher, err := data.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != domain.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := data.GetReservationById(reservationId)
	if err != nil || reservation.Status == domain.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == domain.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
		return nil, errors.New("只能反馈本人开设的咨询")
	}
	return reservation, nil
}

// 咨询师提交反馈
func (tl *TeacherLogic) SubmitFeedbackByTeacher(reservationId string, teacherFullname string, teacherId string,
	studentName string, problem string, solution string, adviceToCenter string, username string,
	userType domain.UserType) (*domain.Reservation, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != domain.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(teacherFullname, "") {
		return nil, errors.New("咨询师姓名为空")
	} else if strings.EqualFold(teacherId, "") {
		return nil, errors.New("咨询师工作证号为空")
	} else if strings.EqualFold(studentName, "") {
		return nil, errors.New("学生姓名为空")
	} else if strings.EqualFold(problem, "") {
		return nil, errors.New("咨询问题为空")
	} else if strings.EqualFold(solution, "") {
		return nil, errors.New("解决方法为空")
	} else if strings.EqualFold(adviceToCenter, "") {
		return nil, errors.New("工作建议为空")
	}
	teacher, err := data.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != domain.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := data.GetReservationById(reservationId)
	if err != nil || reservation.Status == domain.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == domain.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
		return nil, errors.New("只能反馈本人开设的咨询")
	}
	if reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty() {
		sms.SendFeedbackSMS(reservation)
	}
	reservation.TeacherFeedback = domain.TeacherFeedback{
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherId,
		StudentFullname: studentName,
		Problem:         problem,
		Solution:        solution,
		AdviceToCenter:  adviceToCenter,
	}
	if err = data.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 咨询师查看学生信息
func (tl *TeacherLogic) GetStudentInfoByTeacher(reservationId string, username string, userType domain.UserType) (domain.StudentInfo, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	} else if userType != domain.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	}
	teacher, err := data.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != domain.TEACHER {
		return nil, errors.New("权限不足")
	}
	reservation, err := data.GetReservationById(reservationId)
	if err != nil || reservation.Status == domain.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == domain.AVAILABLE {
		return nil, errors.New("咨询未被预约,无法查看")
	} else if !strings.EqualFold(reservation.TeacherUsername, teacher.Username) {
		return nil, errors.New("只能查看本人开设的咨询")
	}
	return reservation.StudentInfo, nil
}
