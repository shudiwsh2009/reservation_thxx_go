package main

import (
	"flag"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"gopkg.in/mgo.v2"
	"time"
)

func main() {
	utils.APP_ENV = *flag.String("app-env", "STAGING", "app environment")
	utils.SMS_UID = *flag.String("sms-uid", "", "sms uid")
	utils.SMS_KEY = *flag.String("sms-key", "", "sms key")
	// 数据库连接
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		fmt.Errorf("连接数据库失败：%v", err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("appointment")
	// 时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Errorf("初始化时区失败：%v", err)
		return
	}
	// Reminder
	now := time.Now().In(location)
	year, month, day := now.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, location)
	from := today.AddDate(0, 0, 1)
	to := today.AddDate(0, 0, 2)
	reservations, err := models.GetReservationsBetweenTime(from, to)
	if err != nil {
		fmt.Errorf("获取咨询列表失败：%v", err)
		return
	}
	for _, reservation := range reservations {
		if reservation.Status == models.RESERVATED {
			utils.SendReminderSMS(reservation)
		}
	}
}
