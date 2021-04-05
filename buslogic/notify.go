package buslogic

import (
	"bytes"
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	SmsSuccessStudent               = "%s你好，你已成功预约星期%s（%d月%d日）%s-%s咨询，地点：%s。电话：62792453。"
	SmsSuccessStudentWithLocation   = "%s你好，你已成功预约星期%s（%d月%d日）%s-%s咨询，地点：%s，咨询方式：%s。电话：62792453。"
	SmsEnSuccessStudent             = "Dear %s, you have successfully made an appointment of advising service for %s (%s %d) from %s to %s in %s. Tel: 62792453."
	SmsEnSuccessStudentWithLocation = "Dear %s, you have successfully made an appointment of advising service for %s (%s %d) from %s to %s in %s. The service will be %s. Tel: 62792453."
	SmsSuccessTeacher               = "%s您好，%s已预约您星期%s（%d月%d日）%s-%s咨询，地点：%s。电话：62792453。"
	SmsSuccessTeacherWithLocation   = "%s您好，%s已预约您星期%s（%d月%d日）%s-%s咨询，地点：%s，咨询方式：%s。电话：62792453。"
	SmsReminderStudent              = "温馨提示：%s你好，你已成功预约明天%s-%s咨询，地点：%s。电话：62792453。"
	SmsReminderStudentWithLocation  = "温馨提示：%s你好，你已成功预约明天%s-%s咨询，地点：%s，咨询方式：%s。电话：62792453。"
	SmsReminderTeacher              = "温馨提示：%s您好，%s已预约您明天%s-%s咨询，地点：%s。电话：62792453。"
	SmsReminderTeacherWithLocation  = "温馨提示：%s您好，%s已预约您明天%s-%s咨询，地点：%s，咨询方式：%s。电话：62792453。"
	SmsFeedbackStudent              = "温馨提示：%s你好，感谢使用我们的一对一咨询服务，请再次登录乐学预约界面，为咨询师反馈评分，帮助我们成长。"
	SmsFeedbackStudentNight         = "温馨提示：%s你好，你的咨询已经结束，请及时通过咨询预约平台填写反馈；如已填过反馈，请忽略本短信提醒。中心电话：62792453。"
	SmsCancelTeacher                = "【预约取消通知】%s咨询师您好，您%d月%d日%s-%s的咨询预约已被取消，请知悉。"
	SmsCancelStudent                = "【预约取消通知】%s同学您好，您%d月%d日%s-%s的咨询因故被取消，请知悉。电话：62792453。"
)

var (
	SMS_ERROR_MSG = map[int]string{
		-1:  "没有该用户账户",
		-2:  "接口密钥不正确",
		-21: "MD5接口密钥加密不正确",
		-3:  "短信数量不足",
		-11: "该用户被禁用",
		-14: "短信内容出现非法字符",
		-4:  "手机号格式不正确",
		-41: "手机号码为空",
		-42: "短信内容为空",
		-51: "短信签名格式不正确",
		-52: "短信签名太长",
		-6:  "IP限制",
	}
)

func (w *Workflow) NotifyMakeReservationSuccess(reservation *model.Reservation, teacher *model.Teacher) error {
	studentSMS := fmt.Sprintf(SmsSuccessStudent, reservation.StudentInfo.Fullname, utils.ChineseShortWeekday[reservation.StartTime.Weekday()],
		reservation.StartTime.Month(), reservation.StartTime.Day(), reservation.StartTime.Format("15:04"),
		reservation.EndTime.Format("15:04"), reservation.TeacherAddress)
	if locationDesc, ok := model.LocationDesc[reservation.StudentInfo.Location]; ok {
		studentSMS = fmt.Sprintf(SmsSuccessStudentWithLocation, reservation.StudentInfo.Fullname, utils.ChineseShortWeekday[reservation.StartTime.Weekday()],
			reservation.StartTime.Month(), reservation.StartTime.Day(), reservation.StartTime.Format("15:04"),
			reservation.EndTime.Format("15:04"), reservation.TeacherAddress, locationDesc)
	}
	if teacher.SmsSuffix != "" {
		studentSMS += teacher.SmsSuffix
	}
	if err := w.sendSMS(reservation.StudentInfo.Mobile, studentSMS); err != nil {
		return err
	}
	go SendEmail("咨询预约成功", studentSMS, []string{reservation.StudentInfo.Email})

	if reservation.InternationalType == model.InternationalTypeChinglish {
		studentSmsEn := fmt.Sprintf(SmsEnSuccessStudent, reservation.StudentInfo.Fullname, utils.EnglishWeekday[reservation.StartTime.Weekday()],
			utils.EnglishShortMonth[reservation.StartTime.Month()], reservation.StartTime.Day(), reservation.StartTime.Format("15:04"),
			reservation.EndTime.Format("15:04"), reservation.TeacherAddressEn)
		if locationDesc, ok := model.LocationDescEn[reservation.StudentInfo.Location]; ok {
			studentSmsEn = fmt.Sprintf(SmsEnSuccessStudentWithLocation, reservation.StudentInfo.Fullname, utils.EnglishWeekday[reservation.StartTime.Weekday()],
				utils.EnglishShortMonth[reservation.StartTime.Month()], reservation.StartTime.Day(), reservation.StartTime.Format("15:04"),
				reservation.EndTime.Format("15:04"), reservation.TeacherAddressEn, locationDesc)
		}
		if teacher.SmsSuffixEn != "" {
			studentSmsEn += teacher.SmsSuffixEn
		}
		if err := w.sendSMS(reservation.StudentInfo.Mobile, studentSmsEn); err != nil {
			return err
		}
		go SendEmail("Appointment of advising service", studentSmsEn, []string{reservation.StudentInfo.Email})
	}

	teacherSMS := fmt.Sprintf(SmsSuccessTeacher, reservation.TeacherFullname, reservation.StudentInfo.Fullname,
		utils.ChineseShortWeekday[reservation.StartTime.Weekday()], reservation.StartTime.Month(), reservation.StartTime.Day(),
		reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"), reservation.TeacherAddress)
	if locationDesc, ok := model.LocationDesc[reservation.StudentInfo.Location]; ok {
		teacherSMS = fmt.Sprintf(SmsSuccessTeacherWithLocation, reservation.TeacherFullname, reservation.StudentInfo.Fullname,
			utils.ChineseShortWeekday[reservation.StartTime.Weekday()], reservation.StartTime.Month(), reservation.StartTime.Day(),
			reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"), reservation.TeacherAddress, locationDesc)
	}
	if err := w.sendSMS(reservation.TeacherMobile, teacherSMS); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) SendReminderSMS(reservation *model.Reservation) error {
	studentSMS := fmt.Sprintf(SmsReminderStudent, reservation.StudentInfo.Fullname, reservation.StartTime.Format("15:04"),
		reservation.EndTime.Format("15:04"), reservation.TeacherAddress)
	if locationDesc, ok := model.LocationDesc[reservation.StudentInfo.Location]; ok {
		studentSMS = fmt.Sprintf(SmsReminderStudentWithLocation, reservation.StudentInfo.Fullname, reservation.StartTime.Format("15:04"),
			reservation.EndTime.Format("15:04"), reservation.TeacherAddress, locationDesc)
	}
	if err := w.sendSMS(reservation.StudentInfo.Mobile, studentSMS); err != nil {
		return err
	}

	teacherSMS := fmt.Sprintf(SmsReminderTeacher, reservation.TeacherFullname, reservation.StudentInfo.Fullname,
		reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"), reservation.TeacherAddress)
	if locationDesc, ok := model.LocationDesc[reservation.StudentInfo.Location]; ok {
		teacherSMS = fmt.Sprintf(SmsReminderTeacherWithLocation, reservation.TeacherFullname, reservation.StudentInfo.Fullname,
			reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"), reservation.TeacherAddress, locationDesc)
	}
	if err := w.sendSMS(reservation.TeacherMobile, teacherSMS); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) SendFeedbackSMS(reservation *model.Reservation) error {
	studentSMS := fmt.Sprintf(SmsFeedbackStudent, reservation.StudentInfo.Fullname)
	if err := w.sendSMS(reservation.StudentInfo.Mobile, studentSMS); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) SendFeedbackNightSMS(reservation *model.Reservation) error {
	studentSMS := fmt.Sprintf(SmsFeedbackStudentNight, reservation.StudentInfo.Fullname)
	if err := w.sendSMS(reservation.StudentInfo.Mobile, studentSMS); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) SendCancelSMS(reservation *model.Reservation, studentFullname, studentMobile string) error {
	studentSms := fmt.Sprintf(SmsCancelStudent, studentFullname, reservation.StartTime.Month(),
		reservation.StartTime.Day(), reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"))
	if err := w.sendSMS(studentMobile, studentSms); err != nil {
		return err
	}
	teacherSMS := fmt.Sprintf(SmsCancelTeacher, reservation.TeacherFullname, reservation.StartTime.Month(), reservation.StartTime.Day(),
		reservation.StartTime.Format("15:04"), reservation.EndTime.Format("15:04"))
	if err := w.sendSMS(reservation.TeacherMobile, teacherSMS); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) sendSMS(mobile string, content string) error {
	if !utils.IsMobile(mobile) {
		return re.NewRErrorCode("手机号格式不正确", nil, re.ErrorFormatMobile)
	}
	if config.Instance().IsStagingEnv() {
		log.Infof("SMOCK Send SMS: \"%s\" to %s", content, mobile)
		return nil
	}
	requestUrl := "http://utf8.sms.webchinese.cn"
	payload := url.Values{
		"Uid":     {config.Instance().SMSUid},
		"Key":     {config.Instance().SMSKey},
		"smsMob":  {mobile},
		"smsText": {content},
	}
	requestBody := bytes.NewBufferString(payload.Encode())
	response, err := http.Post(requestUrl, "application/x-www-form-urlencoded;charset=utf8", requestBody)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	errCode, err := strconv.Atoi(string(responseBody))
	if err != nil {
		return re.NewRError(fmt.Sprintf("短信发送失败：failed to strconv.Atoi %s", responseBody), err)
	}
	if errCode > 0 {
		log.Infof("Send SMS \"%s\" to %s: return %d", content, mobile, errCode)
		return nil
	}
	errMsg, _ := SMS_ERROR_MSG[errCode]
	EmailWarn("短信发送失败", fmt.Sprintf("Fail to send SMS \"%s\" to %s: errCode = %d, errMsg = %s", content, mobile, errCode, errMsg))
	return re.NewRError(fmt.Sprintf("Fail to send SMS \"%s\" to %s: errCode = %d, errMsg = %s", content, mobile, errCode, errMsg), nil)
}

// external 每天20:00发送第二天预约咨询的提醒短信
func (w *Workflow) SendTomorrowReservationReminderSMS() {
	today := utils.BeginOfDay(time.Now())
	from := today.AddDate(0, 0, 1)
	to := today.AddDate(0, 0, 2)
	reservations, err := w.MongoClient().GetReservationsBetweenTime(from, to)
	if err != nil {
		log.Errorf("获取咨询列表失败：%v", err)
		return
	}
	succCnt, failCnt := 0, 0
	for _, reservation := range reservations {
		if reservation.Status == model.ReservationStatusReservated {
			if err = w.SendReminderSMS(reservation); err == nil {
				succCnt++
			} else {
				log.Errorf("发送短信失败：%+v %+v", reservation, err)
				failCnt++
			}
		}
	}
	log.Infof("发送%d个预约记录的提醒短信，成功%d个，失败%d个", succCnt+failCnt, succCnt, failCnt)
}

// external 每天22:00提醒当天咨询的同学填写反馈
func (w *Workflow) SendTodayFeedbackReminderSMS() {
	from := utils.BeginOfDay(time.Now())
	to := from.AddDate(0, 0, 1)
	reservations, err := w.MongoClient().GetReservationsBetweenTime(from, to)
	if err != nil {
		log.Errorf("获取咨询列表失败: %+v", err)
		return
	}
	succCnt, failCnt := 0, 0
	for _, reservation := range reservations {
		if reservation.Status == model.ReservationStatusReservated ||
			reservation.Status == model.ReservationStatusFeedback {
			if err = w.SendFeedbackNightSMS(reservation); err == nil {
				succCnt++
			} else {
				log.Errorf("发送短信失败：%+v %+v", reservation, err)
				failCnt++
			}
		}
	}
	log.Infof("发送%d个预约记录的反馈短信，成功%d个，失败%d个", succCnt+failCnt, succCnt, failCnt)
}
