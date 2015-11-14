package main

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/data"
	"gopkg.in/mgo.v2"
)

func main() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		fmt.Errorf("连接数据库失败:%v", err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	data.Mongo = session.DB("appointment")

	reservations, _ := data.GetReservationsByStudentId("2013213474")
	for _, r := range reservations {
		fmt.Println(r)
	}
}
