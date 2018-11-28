package utils

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"regexp"
	"time"
)

/**
 * 正则：手机号（精确）
 * <p>移动：134(0-8)、135、136、137、138、139、147、150、151、152、157、158、159、178、182、183、184、187、188、198</p>
 * <p>联通：130、131、132、145、155、156、175、176、185、186、166</p>
 * <p>电信：133、153、173、177、180、181、189、199</p>
 * <p>全球星：1349</p>
 * <p>虚拟运营商：170</p>
 */
const (
	RegexMobileExact = "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
)

func IsMobile(mobile string) bool {
	if m, _ := regexp.MatchString(RegexMobileExact, mobile); !m {
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
