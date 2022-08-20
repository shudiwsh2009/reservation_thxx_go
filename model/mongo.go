package model

import (
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	"gopkg.in/mgo.v2"
	"time"
)

var (
	dbTeacher     *mgo.Collection
	dbAdmin       *mgo.Collection
	dbReservation *mgo.Collection
)

type MongoClient struct {
	mongo *mgo.Database
}

func NewMongoClient() *MongoClient {
	var session *mgo.Session
	var err error
	if config.Instance().IsStagingEnv() {
		session, err = mgo.Dial("127.0.0.1:27017")
	} else if config.Instance().IsTsinghuaEnv() {
		session, err = mgo.Dial(config.Instance().MongoHost)
	} else {
		mongoDbDialInfo := mgo.DialInfo{
			Addrs:    []string{config.Instance().MongoHost},
			Timeout:  60 * time.Second,
			Database: config.Instance().MongoAuthDatabase,
			Username: config.Instance().MongoAuthUser,
			Password: config.Instance().MongoAuthPassword,
		}
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	}
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	//defer session.Close()
	session.SetMode(mgo.Eventual, true)
	mongo := session.DB(config.Instance().MongoDatabase)
	dbTeacher = mongo.C("teacher")
	dbAdmin = mongo.C("admin")
	dbReservation = mongo.C("reservation")
	ret := &MongoClient{
		mongo: mongo,
	}
	return ret
}

func (m *MongoClient) EnsureAllIndexes() error {
	var err error

	err = dbTeacher.EnsureIndex(mgo.Index{
		Key: []string{"username", "user_type"},
	})
	if err != nil {
		return err
	}
	err = dbTeacher.EnsureIndex(mgo.Index{
		Key: []string{"fullname", "user_type"},
	})
	if err != nil {
		return err
	}
	err = dbTeacher.EnsureIndex(mgo.Index{
		Key: []string{"mobile", "user_type"},
	})
	if err != nil {
		return err
	}

	err = dbAdmin.EnsureIndex(mgo.Index{
		Key: []string{"username", "user_type"},
	})
	if err != nil {
		return err
	}

	err = dbReservation.EnsureIndex(mgo.Index{
		Key: []string{"student_info.username", "status"},
	})
	if err != nil {
		return err
	}
	err = dbReservation.EnsureIndex(mgo.Index{
		Key: []string{"start_time", "status"},
	})
	if err != nil {
		return err
	}

	return nil
}

// DANGER!!!
func (m *MongoClient) DropAllIndexes() error {
	for _, coll := range []*mgo.Collection{dbTeacher, dbAdmin, dbReservation} {
		err := m.DropIndexes(coll)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MongoClient) DropIndexes(coll *mgo.Collection) error {
	indexes, err := coll.Indexes()
	if err != nil {
		return err
	}
	for _, index := range indexes {
		if index.Name == "_id_" {
			continue
		}
		err = coll.DropIndexName(index.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
