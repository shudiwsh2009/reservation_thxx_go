package utils

import "regexp"

func IsMobile(mobile string) bool {
	if m, _ := regexp.MatchString("(^(13\\d|14[57]|15[^4,\\D]|17[678]|18\\d)\\d{8}|170[059]\\d{7})$", mobile); !m {
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

func IsStudentId(studentId string) bool {
	if m, _ := regexp.MatchString("^\\d{10}$", studentId); !m {
		return false
	}
	return true
}
