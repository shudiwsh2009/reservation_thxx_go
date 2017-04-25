package rerror

import "fmt"

const (
	Ok    = 0
	Check = 1
	// 请求类错误
	ErrorMissingParam = 10
	ErrorInvalidParam = 11
	// 格式类错误
	ErrorFormatMobile          = 51
	ErrorFormatEmail           = 52
	ErrorFormatWeekday         = 53
	ErrorFormatStudentUsername = 54
	// 账户类错误
	ErrorExpireSession                        = 100
	ErrorNoLogin                              = 101
	ErrorNotAuthorized                        = 102
	ERROR_LOGIN_PASSWORD_WRONG                = 103
	ERROR_LOGIN_PWDCHANGE_OLDPWD_MISMATCH     = 104
	ErrorNoUser                               = 105
	ERROR_EXIST_USERNAME                      = 106
	ERROR_LOGIN_PWDCHANGE_OLDPWD_EQUAL_NEWPED = 107
	ERROR_LOGIN_PWDCHANGE_INFO_MISMATCH       = 108
	ERROR_LOGIN_PWDCHANGE_VERIFY_CODE_WRONG   = 110
	ERROR_NO_STUDENT                          = 111
	// 通用逻辑类错误
	ErrorFeedbackAvailableReservation          = 201
	ErrorFeedbackFutureReservation             = 202
	ErrorFeedbackOtherReservation              = 203
	ErrorEditReservationEndTimeBeforeStartTime = 204
	ErrorEditReservationTeacherTimeConflict    = 205
	ErrorEditReservatedReservation             = 206
	ErrorEditOutdatedReservation               = 207
	ErrorViewAvailableReservationStudentInfo   = 208
	ERROR_START_TIME_MISMATCH                  = 209
	ERROR_SEND_SMS                             = 210
	// 管理员错误
	ERROR_ADMIN_SET_RESERVATED_RESERVATION       = 306
	ERROR_ADMIN_ARCHIVE_NUMBER_ALREADY_EXIST     = 308
	ERROR_ADMIN_EXPORT_STUDENT_NO_ARCHIVE_NUMBER = 309
	ERROR_ADMIN_NO_RESERVATIONS_TODAY            = 310
	// 学生错误
	ErrorStudentAlreadyHaveReservation                = 401
	ErrorStudentMakeOutdatedReservation               = 402
	ErrorStudentMakeReservatedReservation             = 403
	ERROR_STUDENT_MAKE_NOT_BINDED_TEACHER_RESERVATION = 404
	ErrorStudentMakeReservationTooEarly               = 405
	// 咨询师错误
	ERROR_TEACHER_VIEW_OTHER_STUDENT = 501
	ErrorTeacherEditOtherReservation = 502
	ErrorTeacherViewOtherReservation = 503
	// 数据库类错误
	ErrorDatabase    = 1000
	ERROR_ID_INVALID = 1001

	ErrorUnknown = -1
)

func ReturnMessage(code int, args ...interface{}) string {
	switch code {
	case Ok:
		return "OK"
	case Check:
		return "CHECK"
	case ErrorMissingParam:
		return fmt.Sprintf("参数缺失：%s", args...)
	case ErrorInvalidParam:
		return fmt.Sprintf("参数格式错误：%s", args...)
	case ErrorExpireSession:
		return "会话过期，请重新登录"
	case ErrorNoLogin:
		return "请先登录"
	case ErrorNotAuthorized:
		return "权限不足"
	case ERROR_LOGIN_PASSWORD_WRONG:
		return "用户名或密码不正确"
	case ERROR_LOGIN_PWDCHANGE_OLDPWD_MISMATCH:
		return "旧登录密码不正确"
	case ErrorNoUser:
		return "未找到用户"
	case ERROR_EXIST_USERNAME:
		return "用户名已被注册"
	case ERROR_LOGIN_PWDCHANGE_OLDPWD_EQUAL_NEWPED:
		return "新密码不能与原有密码一样"
	case ERROR_LOGIN_PWDCHANGE_VERIFY_CODE_WRONG:
		return "验证码错误或已过期"
	case ERROR_LOGIN_PWDCHANGE_INFO_MISMATCH:
		return "信息不匹配"
	case ERROR_NO_STUDENT:
		return "学生未注册"
	case ErrorFormatMobile:
		return "手机号格式不正确"
	case ErrorFormatEmail:
		return "邮箱格式不正确"
	case ErrorFormatWeekday:
		return "星期格式不正确"
	case ErrorFormatStudentUsername:
		return "学号格式不正确"
	case ERROR_SEND_SMS:
		return "发送短信失败，请稍后重试"
	case ErrorEditReservationEndTimeBeforeStartTime:
		return "开始时间不能晚于结束时间"
	case ErrorEditReservationTeacherTimeConflict:
		return "咨询师时间有冲突"
	case ErrorEditReservatedReservation:
		return "不能编辑已预约的咨询"
	case ErrorEditOutdatedReservation:
		return "不能编辑已过期咨询"
	case ErrorFeedbackAvailableReservation:
		return "不能反馈未被预约的咨询"
	case ErrorFeedbackFutureReservation:
		return "不能反馈还未开始的咨询"
	case ERROR_ADMIN_SET_RESERVATED_RESERVATION:
		return "不能设定已被预约的咨询"
	case ERROR_START_TIME_MISMATCH:
		return "开始时间不匹配"
	case ERROR_ADMIN_ARCHIVE_NUMBER_ALREADY_EXIST:
		return "档案号已存在，请重新分配"
	case ERROR_ADMIN_EXPORT_STUDENT_NO_ARCHIVE_NUMBER:
		return "请先分配档案号"
	case ERROR_ADMIN_NO_RESERVATIONS_TODAY:
		return "今日无咨询"
	case ErrorStudentAlreadyHaveReservation:
		return "你好！你已有一个咨询预约，请完成这次咨询后再预约下一次，或致电62792453取消已有预约。"
	case ErrorStudentMakeOutdatedReservation:
		return "不能预约已过期咨询"
	case ErrorStudentMakeReservatedReservation:
		return "不能预约已被预约的咨询"
	case ERROR_STUDENT_MAKE_NOT_BINDED_TEACHER_RESERVATION:
		return "只能预约匹配咨询师的咨询"
	case ErrorStudentMakeReservationTooEarly:
		return fmt.Sprintf("距咨询开始不足%s，无法预约", args...)
	case ErrorFeedbackOtherReservation:
		return "只能反馈自己预约的咨询"
	case ERROR_TEACHER_VIEW_OTHER_STUDENT:
		return "只能查看本人绑定的学生"
	case ErrorTeacherEditOtherReservation:
		return "只能编辑本人开设的咨询"
	case ErrorTeacherViewOtherReservation:
		return "只能查看本人开设的资讯"
	case ErrorViewAvailableReservationStudentInfo:
		return "咨询未被预约，无法查看"
	case ErrorDatabase:
		return "获取数据失败"
	case ERROR_ID_INVALID:
		return "ID值不合法"
	case ErrorUnknown:
		return "未知错误"
	default:
		return "服务器开小差了，请稍后重试~"
	}
}
