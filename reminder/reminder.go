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
	appEnv := flag.String("app-env", "STAGING", "app environment")
	smsUid := flag.String("sms-uid", "", "sms uid")
	smsKey := flag.String("sms-key", "", "sms key")
	flag.Parse()
	utils.APP_ENV = *appEnv
	utils.SMS_UID = *smsUid
	utils.SMS_KEY = *smsKey
	// 数据库连接
	mongoDbDialInfo := mgo.DialInfo{
		Addrs:		[]string{"127.0.0.1:27017"},
		Timeout:	60 * time.Second,
		Database:	"admin",
		Username:	"admin",
		Password:	"THXXFZZX",
	}
	var session *mgo.Session
	var err error
	if utils.APP_ENV == "ONLINE" {
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	} else {
		session, err = mgo.Dial("127.0.0.1:27017")
	}
	if err != nil {
		fmt.Printf("连接数据库失败：%v\n", err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("appointment")
	// 时区
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("初始化时区失败：%v\n", err)
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
