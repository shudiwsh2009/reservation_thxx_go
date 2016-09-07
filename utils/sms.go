package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	SMS_SUCCESS_STUDENT  = "%s你好，你已成功预约星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼407室。电话：62792453。"
	SMS_SUCCESS_TEACHER  = "%s您好，%s已预约您星期%s（%d月%d日）%s-%s咨询，地点：紫荆C楼407室。电话：62792453。"
	SMS_REMINDER_STUDENT = "温馨提示：%s你好，你已成功预约明天%s-%s咨询，地点：紫荆C楼407室。电话：62792453。"
	SMS_REMINDER_TEACHER = "温馨提示：%s您好，%s已预约您明天%s-%s咨询，地点：紫荆C楼407室。电话：62792453。"
	SMS_FEEDBACK_STUDENT = "温馨提示：%s你好，感谢使用我们的一对一咨询服务，请再次登录乐学预约界面，为咨询师反馈评分，帮助我们成长。"
)

var (
	APP_ENV = ""
	SMS_UID = ""
	SMS_KEY = ""
)

func SendSuccessSMS(reservation *models.Reservation) error {
	studentSMS := fmt.Sprintf(SMS_SUCCESS_STUDENT, reservation.StudentInfo.Name, Weekdays[reservation.StartTime.Weekday()],
		reservation.StartTime.Month(), reservation.StartTime.Day(), reservation.StartTime.Format("15:04"),
		reservation.EndTime.Format("15:04"))
	if err := sendSMS(reservation.StudentInfo.Mobile, studentSMS); err != nil {
		return err
	}
	teacherSMS := fmt.Sprintf(SMS_SUCCESS_TEACHER, reservation.TeacherFullname, reservation.StudentInfo.Name,
		Weekdays[reservation.StartTime.Weekday()], reservation.StartTime.Month(), reservation.StartTime.Day(),
		reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"))
	if err := sendSMS(reservation.TeacherMobile, teacherSMS); err != nil {
		return err
	}
	return nil
}

func SendReminderSMS(reservation *models.Reservation) error {
	studentSMS := fmt.Sprintf(SMS_REMINDER_STUDENT, reservation.StudentInfo.Name, reservation.StartTime.Format("15:04"),
		reservation.EndTime.Format("15:04"))
	if err := sendSMS(reservation.StudentInfo.Mobile, studentSMS); err != nil {
		return err
	}
	teacherSMS := fmt.Sprintf(SMS_REMINDER_TEACHER, reservation.TeacherFullname, reservation.StudentInfo.Name,
		reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"))
	if err := sendSMS(reservation.TeacherMobile, teacherSMS); err != nil {
		return err
	}
	return nil
}

func SendFeedbackSMS(reservation *models.Reservation) error {
	studentSMS := fmt.Sprintf(SMS_FEEDBACK_STUDENT, reservation.StudentInfo.Name)
	if err := sendSMS(reservation.StudentInfo.Mobile, studentSMS); err != nil {
		return err
	}
	return nil
}

func sendSMS(mobile string, content string) error {
	if m := IsMobile(mobile); !m {
		return errors.New("手机号格式不正确")
	}
	if !strings.EqualFold(APP_ENV, "ONLINE") || len(SMS_UID) == 0 || len(SMS_KEY) == 0 {
		fmt.Printf("Send SMS: \"%s\" to %s.\n", content, mobile)
		return nil
	}
	requestUrl := "http://utf8.sms.webchinese.cn"
	payload := url.Values{
		"Uid":     {SMS_UID},
		"Key":     {SMS_KEY},
		"smsMob":  {mobile},
		"smsText": {content},
	}
	requestBody := bytes.NewBufferString(payload.Encode())
	response, err := http.Post(requestUrl, "application/x-www-form-urlencoded;charset=utf8", requestBody)
	if err != nil {
		return errors.New("短信发送失败")
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.New("短信发送失败")
	}
	fmt.Println(string(responseBody))
	if code, err := strconv.Atoi(string(responseBody)); err != nil || code < 0 {
		return errors.New(fmt.Sprintf("短信发送失败,ErrCode:%d", code))
	}
	return nil
}
