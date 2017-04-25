package utils

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"regexp"
	"time"
)

//正则表达式为：
//未转义：(^(13\d|15[^4\D]|17[13678]|18\d)\d{8}|170[^346\D]\d{7})$。
//已转义：(^(13\\d|15[^4\\D]|17[13678]|18\\d)\\d{8}|170[^346\\D]\\d{7})$。
//
//默认 14x 上网卡号段为无效号码，如果希望其为有效号码，则正则表达式为：
//未转义：(^(13\d|14[57]|15[^4\D]|17[13678]|18\d)\d{8}|170[^346\D]\d{7})$。
//已转义：(^(13\\d|14[57]|15[^4\\D]|17[13678]|18\\d)\\d{8}|170[^346\\D]\\d{7})$。
func IsMobile(mobile string) bool {
	if m, _ := regexp.MatchString("(^(13\\d|15[^4\\D]|17[13678]|18\\d)\\d{8}|170[^346\\D]\\d{7})$", mobile); !m {
		return false
	}
	return true
}

func IsEmail(email string) bool {
	if m, _ := regexp.MatchString("^([a-z0-9A-Z]+[-|\\.]?)+[a-z0-9A-Z]@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\\.)+[a-zA-Z]{2,}$",
		email); !m {
		return false
	}
	return true
}

func IsStudentUsername(studentId string) bool {
	if m, _ := regexp.MatchString("^\\d{10}$", studentId); !m {
		return false
	}
	return true
}

func ParseStudentId(studentId string) (string, error) {
	if !IsStudentUsername(studentId) {
		return "", rerror.NewRError(fmt.Sprintf("fail to parse studentId %s", studentId), nil)
	}
	switch studentId[4] {
	case '0':
		return studentId[2:4] + "级", nil
	case '2':
		return studentId[2:4] + "硕", nil
	case '3':
		return studentId[2:4] + "博", nil
	}
	return "", rerror.NewRError(fmt.Sprintf("unknown degree and grade %s", studentId), nil)
}

func StringToWeekday(weekday string) (time.Weekday, error) {
	switch weekday {
	case "Sunday":
		return time.Sunday, nil
	case "Monday":
		return time.Monday, nil
	case "Tuesday":
		return time.Tuesday, nil
	case "Wednesday":
		return time.Wednesday, nil
	case "Thursday":
		return time.Thursday, nil
	case "Friday":
		return time.Friday, nil
	case "Saturday":
		return time.Saturday, nil
	}
	return 0, rerror.NewRErrorCode("星期格式错误", nil, rerror.ErrorFormatWeekday)
}
