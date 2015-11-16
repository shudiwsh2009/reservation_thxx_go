package models

import "gopkg.in/mgo.v2/bson"

type UserType int

const (
	STUDENT UserType = 1 + iota
	TEACHER
	ADMIN
)

var userTypes = [...]string{
	"STUDENT",
	"TEACHER",
	"ADMIN",
}

func (ut UserType) String() string {
	return userTypes[ut-1]
}

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	Fullname string        `bson:"fullname"`
	Mobile   string        `bson:"mobile"`
	UserType UserType      `bson:"userType"`
}
