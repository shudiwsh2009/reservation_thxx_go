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
	if strings.EqualFold(userId, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(startTime, "") {
		return nil, errors.New("开始时间为空")
	} else if strings.EqualFold(endTime, "") {
		return nil, errors.New("结束时间为空")
	} else if strings.EqualFold(teacherFullname, "") {
		return nil, errors.New("咨询师姓名为空")
	} else if strings.EqualFold(teacherMobile, "") {
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
	teacher.Fullname = teacherFullname
	teacher.Mobile = teacherMobile
	if err = models.UpsertUser(teacher); err != nil {
		return nil, errors.New("数据获取失败")
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
	if strings.EqualFold(userId, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
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
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil {
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
	}
	teacher, err := models.GetUserById(userId)
	if err != nil {
		return nil, errors.New("咨询师账户失效")
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if !strings.EqualFold(teacher.Username, reservation.TeacherUsername) {
		return nil, errors.New("只能编辑本人开设的咨询")
	}
	teacher.Fullname = teacherFullname
	teacher.Mobile = teacherMobile
	if err = models.UpsertUser(teacher); err != nil {
		return nil, errors.New("数据获取失败")
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
	if strings.EqualFold(userId, "") {
		return 0, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return 0, errors.New("权限不足")
	} else if reservationIds == nil {
		return 0, errors.New("咨询Id列表为空")
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
	if strings.EqualFold(userId, "") {
		return 0, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return 0, errors.New("权限不足")
	} else if reservationIds == nil {
		return 0, errors.New("咨询Id列表为空")
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
	if strings.EqualFold(userId, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
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
	if strings.EqualFold(userId, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
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
	if strings.EqualFold(userId, "") {
		return nil, errors.New("请先登录")
	} else if userType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else if strings.EqualFold(reservationId, "") {
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
