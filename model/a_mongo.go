package model

import (
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"gopkg.in/mgo.v2"
	"time"
)

var (
	dbLegacyUser *mgo.Collection
	dbLegacyReservation *mgo.Collection
)

type LegacyMongoClient struct {
	mongo *mgo.Database
}

func NewLegacyMongoClient() *LegacyMongoClient {
	var session *mgo.Session
	var err error
	if config.Instance().IsStagingEnv() {
		session, err = mgo.Dial("127.0.0.1:27017")
	} else {
		mongoDbDialInfo := mgo.DialInfo{
			Addrs:        []string{"127.0.0.1:27017"},
			Timeout:    60 * time.Second,
			Database:    "admin",
			Username:    "admin",
			Password:    "THXXFZZX",
		}
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	}
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	//defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	mongo := session.DB("appointment")
	dbLegacyUser = mongo.C("user")
	dbLegacyReservation = mongo.C("appointment")
	ret := &LegacyMongoClient{
		mongo: mongo,
	}
	return ret
}
