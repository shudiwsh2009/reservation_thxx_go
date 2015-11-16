package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"strings"
	"time"
)

type StudentLogic struct {
}

// 学生预约咨询
func (sl *StudentLogic) MakeReservationByStudent(reservationId string, name string, gender string,
	studentId string, school string, hometown string, mobile string, email string, experience string,
	problem string) (*models.Reservation, error) {
	if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(name, "") {
		return nil, errors.New("姓名为空")
	} else if strings.EqualFold(gender, "") {
		return nil, errors.New("性别为空")
	} else if strings.EqualFold(studentId, "") {
		return nil, errors.New("学号为空")
	} else if strings.EqualFold(school, "") {
		return nil, errors.New("院系为空")
	} else if strings.EqualFold(hometown, "") {
		return nil, errors.New("生源地为空")
	} else if strings.EqualFold(mobile, "") {
		return nil, errors.New("手机号为空")
	} else if strings.EqualFold(email, "") {
		return nil, errors.New("邮箱为空")
	} else if strings.EqualFold(experience, "") {
		return nil, errors.New("咨询经历为空")
	} else if strings.EqualFold(problem, "") {
		return nil, errors.New("咨询问题为空")
	} else if !utils.IsStudentId(studentId) {
		return nil, errors.New("学号不正确")
	} else if !utils.IsMobile(mobile) {
		return nil, errors.New("手机号格式不正确")
	} else if !utils.IsEmail(email) {
		return nil, errors.New("邮箱格式不正确")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.Before(time.Now().Local()) {
		return nil, errors.New("咨询已过期")
	} else if reservation.Status != models.AVAILABLE {
		return nil, errors.New("咨询已被预约")
	}
	studentReservations, err := models.GetReservationsByStudentId(studentId)
	if err != nil {
		return nil, errors.New("数据获取失败")
	}
	for _, r := range studentReservations {
		if r.Status == models.RESERVATED && r.StartTime.After(time.Now().Local()) {
			return nil, errors.New("你好！你已有一个咨询预约，请完成这次咨询后再预约下一次，或致电62792453取消已有预约。")
		}
	}
	reservation.StudentInfo = models.StudentInfo{
		Name:       name,
		Gender:     gender,
		StudentId:  studentId,
		School:     school,
		Hometown:   hometown,
		Mobile:     mobile,
		Email:      email,
		Experience: experience,
		Problem:    problem,
	}
	reservation.Status = models.RESERVATED
	err = models.UpsertReservation(reservation)
	if err != nil {
		return nil, errors.New("获取数据失败")
	}

	// send success sms
	if checkReservation, err := models.GetReservationById(reservationId); err == nil &&
		checkReservation.Status == models.RESERVATED && strings.EqualFold(checkReservation.StudentInfo.Mobile, mobile) {
		utils.SendSuccessSMS(checkReservation)
	}
	return reservation, nil
}

// 学生拉取反馈
func (sl *StudentLogic) GetFeedbackByStudent(reservationId string, studentId string) (*models.Reservation, error) {
	if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(studentId, "") || !utils.IsStudentId(studentId) {
		return nil, errors.New("学号不正确")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.StudentInfo.StudentId, studentId) {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	return reservation, nil
}

// 学生反馈
func (sl *StudentLogic) SubmitFeedbackByStudent(reservationId string, name string, problem string, choices string,
	score string, feedback string, studentId string) (*models.Reservation, error) {
	if strings.EqualFold(reservationId, "") {
		return nil, errors.New("咨询已下架")
	} else if strings.EqualFold(name, "") {
		return nil, errors.New("姓名为空")
	} else if strings.EqualFold(problem, "") {
		return nil, errors.New("咨询问题为空")
	} else if strings.EqualFold(choices, "") {
		return nil, errors.New("选项为空")
	} else if strings.EqualFold(score, "") {
		return nil, errors.New("总评为空")
	} else if strings.EqualFold(feedback, "") {
		return nil, errors.New("反馈为空")
	} else if strings.EqualFold(studentId, "") || !utils.IsStudentId(studentId) {
		return nil, errors.New("学号不正确")
	}
	reservation, err := models.GetReservationById(reservationId)
	if err != nil {
		return nil, errors.New("咨询已下架")
	} else if reservation.StartTime.After(time.Now().Local()) {
		return nil, errors.New("咨询未开始,暂不能反馈")
	} else if reservation.Status == models.AVAILABLE {
		return nil, errors.New("咨询未被预约,不能反馈")
	} else if !strings.EqualFold(reservation.StudentInfo.StudentId, studentId) {
		return nil, errors.New("只能反馈本人预约的咨询")
	}
	reservation.StudentFeedback = models.StudentFeedback{
		Name:     name,
		Problem:  problem,
		Choices:  choices,
		Score:    score,
		Feedback: feedback,
	}
	if err = models.UpsertReservation(reservation); err != nil {
		return nil, errors.New("数据获取失败")
	}
	return reservation, nil
}
