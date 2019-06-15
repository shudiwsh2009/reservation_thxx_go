package main

import (
	"flag"
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
)

func main() {
	conf := flag.String("conf", "deploy/thxx.conf", "conf file path")
	isStaging := flag.Bool("staging", true, "is staging server")
	method := flag.String("method", "", "method")
	username := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	userType := flag.Int("usertype", model.UserTypeUnknown, "usertype")
	days := flag.Int("days", 0, "shift reservation time in days")
	flag.Parse()

	config.InitWithParams(*conf, *isStaging)
	log.Infof("config loaded: %+v", *config.Instance())
	workflow := buslogic.NewWorkflow()

	if *method == "reminder" {
		// 每晚发送第二天咨询的提醒短信
		workflow.SendTomorrowReservationReminderSMS()
	} else if *method == "feedback-reminder" {
		workflow.SendTodayFeedbackReminderSMS()
	} else if *method == "reset-user-password" {
		// 重置用户密码
		if err := workflow.ResetUserPassword(*username, *userType, *password); err != nil {
			log.Errorf("reset user password failed, err: %+v", err)
			return
		}
	} else if *method == "shift-reservation-time" {
		// 将所有咨询的开始时间和结束时间变更数天
		if err := workflow.ShiftReservationTimeInDays(*days); err != nil {
			log.Errorf("fail to shift reservation: %+v", err)
			return
		}
	} else if *method == "add-new-admin" {
		// 添加新管理员
		if _, err := workflow.AddNewAdmin(*username, *password); err != nil {
			log.Errorf("fail to add new admin: %+v", err)
			return
		}
	}
	log.Info("Success")
}
