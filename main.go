package main

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/data"
	"github.com/shudiwsh2009/reservation_thxx_go/domain"
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

	reservationIds := []string{"5646f0f7a56d4188cf603efb", "5646f0eaa56d4188cf603efa"}
	var reservations []*domain.Reservation
	for _, reservationId := range reservationIds {
		reservation, err := data.GetReservationById(reservationId)
		if err != nil {
			continue
		}
		reservations = append(reservations, reservation)
	}
	for _, reservation := range reservations {
		fmt.Println(reservation)
	}

}
