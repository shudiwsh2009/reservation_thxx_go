package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"strings"
	"time"
)

type AdminLogic struct {
}

// 管理员添加咨询
func (al *AdminLogic) AddReservationByAdmin(startTime string, endTime string, teacherUsername string,
	teacherFullname string, teacherMobile string, userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(startTime) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endTime) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
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
	teacher, err := models.GetUserByUsername(teacherUsername)
	if err != nil {
		if teacher, err = models.AddFullUser(teacherUsername, TeacherDefaultPassword, teacherFullname,
			teacherMobile, models.TEACHER); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertUser(teacher); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation, err := models.AddReservation(start, end, teacher.Fullname, teacher.Username, teacher.Mobile)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员编辑咨询
func (al *AdminLogic) EditReservationByAdmin(reservationId string, startTime string, endTime string,
	teacherUsername string, teacherFullname string, teacherMobile string, userId string,
	userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(startTime) == 0 {
		return nil, errors.New("开始时间为空")
	} else if len(endTime) == 0 {
		return nil, errors.New("结束时间为空")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工号为空")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherMobile) == 0 {
		return nil, errors.New("咨询师手机号为空")
	} else if !utils.IsMobile(teacherMobile) {
		return nil, errors.New("咨询师手机号格式不正确")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
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
	teacher, err := models.GetUserByUsername(teacherUsername)
	if err != nil {
		if teacher, err = models.AddFullUser(teacherUsername, TeacherDefaultPassword, teacherFullname,
			teacherMobile, models.TEACHER); err != nil {
			return nil, errors.New("获取数据失败")
		}
	} else if teacher.UserType != models.TEACHER {
		return nil, errors.New("权限不足")
	} else {
		teacher.Fullname = teacherFullname
		teacher.Mobile = teacherMobile
		if err = models.UpsertUser(teacher); err != nil {
			return nil, errors.New("获取数据失败")
		}
	}
	reservation.StartTime = start
	reservation.EndTime = end
	reservation.TeacherUsername = teacher.Username
	reservation.TeacherFullname = teacher.Fullname
	reservation.TeacherMobile = teacher.Mobile
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员删除咨询
func (al *AdminLogic) RemoveReservationsByAdmin(reservationIds []string, userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for _, reservationId := range reservationIds {
		if reservation, err := models.GetReservationById(reservationId); err == nil {
			reservation.Status = models.DELETED
			models.UpsertReservation(reservation)
			removed++
		}
	}
	return removed, nil
}

// 管理员取消预约
func (al *AdminLogic) CancelReservationsByAdmin(reservationIds []string, userId string, userType models.UserType) (int, error) {
	if len(userId) == 0 {
		return 0, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return 0, errors.New("权限不足")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return 0, errors.New("管理员账户出错,请联系技术支持")
	}
	removed := 0
	for _, reservationId := range reservationIds {
		reseravtion, err := models.GetReservationById(reservationId)
		if err != nil || reseravtion.Status == models.DELETED {
			continue
		}
		if reseravtion.Status == models.RESERVATED && reseravtion.StartTime.After(time.Now().Local()) {
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

// 管理员拉取反馈
func (al *AdminLogic) GetFeedbackByAdmin(reservationId string, userId string, userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	}
	return reservation, nil
}

// 管理员提交反馈
func (al *AdminLogic) SubmitFeedbackByAdmin(reservationId string, teacherFullname string, teacherUsername string,
	studentFullname string, problem string, solution string, adviceToCenter string, userId string,
	userType models.UserType) (*models.Reservation, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	} else if len(teacherFullname) == 0 {
		return nil, errors.New("咨询师姓名为空")
	} else if len(teacherUsername) == 0 {
		return nil, errors.New("咨询师工作证号为空")
	} else if len(studentFullname) == 0 {
		return nil, errors.New("学生姓名为空")
	} else if len(problem) == 0 {
		return nil, errors.New("咨询问题为空")
	} else if len(solution) == 0 {
		return nil, errors.New("解决方法为空")
	} else if len(adviceToCenter) == 0 {
		return nil, errors.New("工作建议为空")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(teacherUsername, reservation.TeacherUsername) {
		return nil, errors.New("咨询师工号不匹配")
	}
	if reservation.TeacherFeedback.IsEmpty() && reservation.StudentFeedback.IsEmpty() {
		utils.SendFeedbackSMS(reservation)
	}
	reservation.TeacherFeedback = models.TeacherFeedback{
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherUsername,
		StudentFullname: studentFullname,
		Problem:         problem,
		Solution:        solution,
		AdviceToCenter:  adviceToCenter,
	}
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}

// 管理员查看学生信息
func (al *AdminLogic) GetStudentInfoByAdmin(reservationId string, userId string, userType models.UserType) (*models.StudentInfo, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	} else if len(reservationId) == 0 {
		return nil, errors.New("咨询已下架")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil || reservation.Status == models.DELETED {
		return nil, errors.New("咨询已下架")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,无法查看")
	}
	return &reservation.StudentInfo, nil
}

// 管理员导出咨询
func (al *AdminLogic) ExportReservationsByAdmin(reservationIds []string, userId string, userType models.UserType) (string, error) {
	if len(userId) == 0 {
		return "", errors.New("请先登录")
	} else if userType != models.ADMIN {
		return "", errors.New("权限不足")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return "", errors.New("管理员账户出错,请联系技术支持")
	}
	var reservations []*models.Reservation
	for _, reservationId := range reservationIds {
		reservation, err := models.GetReservationById(reservationId)
		if err != nil {
			continue
		}
		reservations = append(reservations, reservation)
	}
	filename := "export_" + time.Now().In(utils.Location).Format(utils.DATE_PATTERN) + utils.ExcelSuffix
	if len(reservations) == 0 {
		return "", nil
	}
	if err = utils.ExportReservationsToExcel(reservations, filename); err != nil {
		return "", err
	}
	return "/" + utils.ExportFolder + filename, nil
}

// 查找咨询师
// 查找顺序:全名 > 工号 > 手机号
func (al *AdminLogic) SearchTeacherByAdmin(teacherFullname string, teacherUsername string, teacherMobile string, userId string, userType models.UserType) (*models.User, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	} else if userType != models.ADMIN {
		return nil, errors.New("权限不足")
	}
	admin, err := models.GetUserById(userId)
	if err != nil || admin.UserType != models.ADMIN {
		return nil, errors.New("管理员账户出错,请联系技术支持")
	}
	if !strings.EqualFold(teacherFullname, "") {
		user, err := models.GetUserByFullname(teacherFullname)
		if err == nil {
			return user, nil
		}
	}
	if !strings.EqualFold(teacherUsername, "") {
		user, err := models.GetUserByUsername(teacherUsername)
		if err == nil {
			return user, nil
		}
	}
	if !strings.EqualFold(teacherMobile, "") {
		user, err := models.GetUserByMobile(teacherMobile)
		if err == nil {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}
