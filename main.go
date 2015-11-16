package main

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/data"
	"gopkg.in/mgo.v2"
	"time"
	"github.com/shudiwsh2009/reservation_thxx_go/util"
	"os"
)

func init() {
	// 初始化数据库连接
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		fmt.Errorf("连接数据库失败：%v", err)
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	data.Mongo = session.DB("appointment")
	// 初始化时区
	if util.Location, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		fmt.Errorf("初始化时区失败：%v", err)
		os.Exit(1)
	}
}

func main() {

//	reservationIds := []string{"5646f0f7a56d4188cf603efb", "5646f0eaa56d4188cf603efa"}
//	var reservations []*domain.Reservation
//	for _, reservationId := range reservationIds {
//		reservation, err := data.GetReservationById(reservationId)
//		if err != nil {
//			continue
//		}
//		reservations = append(reservations, reservation)
//	}
//	for _, reservation := range reservations {
//		fmt.Println(reservation)
//	}
	t := time.Now().In(util.Location)
	fmt.Println(t.Format("2006-01-02 15:04"))

}
