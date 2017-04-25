package models

import "gopkg.in/mgo.v2/bson"

type UserType int

const (
	UNKNWONUSER UserType = iota
	STUDENT
	TEACHER
	ADMIN
)

var userTypes = [...]string{
	"UNKNOWN",
	"STUDENT",
	"TEACHER",
	"ADMIN",
}

func (ut UserType) String() string {
	return userTypes[ut]
}

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	Fullname string        `bson:"fullname"`
	Mobile   string        `bson:"mobile"`
	UserType UserType      `bson:"userType_GO"`
	Gender   string        `bson:"gender"`
	Major    string        `bson:"major"`
	Academic string        `bson:"academic"`
	Aptitude string        `bson:"aptitude"`
	Problem  string        `bson:"problem"`
	Address  string        `bson:"address"` // 咨询师地址
}
