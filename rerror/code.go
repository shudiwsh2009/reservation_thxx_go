package rerror

import "fmt"

const (
	Ok = 0
	// 请求类错误
	ErrorMissingParam = 10
	ErrorInvalidParam = 11
	// 格式类错误
	ErrorFormatMobile          = 51
	ErrorFormatEmail           = 52
	ErrorFormatWeekday         = 53
	ErrorFormatStudentUsername = 54
	// 账户类错误
	ErrorExpireSession      = 100
	ErrorNoLogin            = 101
	ErrorNotAuthorized      = 102
	ErrorLoginPasswordWrong = 103
	ErrorNoUser             = 105
	// 通用逻辑类错误
	ErrorFeedbackAvailableReservation          = 201
	ErrorFeedbackFutureReservation             = 202
	ErrorFeedbackOtherReservation              = 203
	ErrorEditReservationEndTimeBeforeStartTime = 204
	ErrorEditReservationTeacherTimeConflict    = 205
	ErrorEditReservatedReservation             = 206
	ErrorEditOutdatedReservation               = 207
	ErrorViewAvailableReservationStudentInfo   = 208
	// 管理员错误
	ErrorAdminNoExportableReservations = 310
	ErrorAdminExportReservationFailure = 320
	// 学生错误
	ErrorStudentAlreadyHaveReservation    = 401
	ErrorStudentMakeOutdatedReservation   = 402
	ErrorStudentMakeReservatedReservation = 403
	ErrorStudentMakeReservationTooEarly   = 405
	// 咨询师错误
	ErrorTeacherEditOtherReservation = 502
	ErrorTeacherViewOtherReservation = 503
	// 数据库类错误
	ErrorDatabase = 1000

	ErrorUnknown = -1
)

func ReturnMessage(code int, args ...interface{}) string {
	switch code {
	case Ok:
		return "OK"
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
	case ErrorLoginPasswordWrong:
		return "用户名或密码不正确"
	case ErrorNoUser:
		return "未找到用户"
	case ErrorFormatMobile:
		return "手机号格式不正确"
	case ErrorFormatEmail:
		return "邮箱格式不正确"
	case ErrorFormatWeekday:
		return "星期格式不正确"
	case ErrorFormatStudentUsername:
		return "学号格式不正确"
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
	case ErrorAdminNoExportableReservations:
		return "没有可以导出的咨询"
	case ErrorAdminExportReservationFailure:
		return "导出咨询失败，请联系技术人员"
	case ErrorStudentAlreadyHaveReservation:
		return "你好！你已有一个咨询预约，请完成这次咨询后再预约下一次，或致电62792453取消已有预约。"
	case ErrorStudentMakeOutdatedReservation:
		return "不能预约已过期咨询"
	case ErrorStudentMakeReservatedReservation:
		return "不能预约已被预约的咨询"
	case ErrorStudentMakeReservationTooEarly:
		return fmt.Sprintf("距咨询开始不足%s，无法预约", args...)
	case ErrorFeedbackOtherReservation:
		return "只能反馈自己预约的咨询"
	case ErrorTeacherEditOtherReservation:
		return "只能编辑本人开设的咨询"
	case ErrorTeacherViewOtherReservation:
		return "只能查看本人开设的资讯"
	case ErrorViewAvailableReservationStudentInfo:
		return "咨询未被预约，无法查看"
	case ErrorDatabase:
		return "获取数据失败"
	case ErrorUnknown:
		return "未知错误"
	default:
		return "服务器开小差了，请稍后重试~"
	}
}
