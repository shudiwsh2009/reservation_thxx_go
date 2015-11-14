package domain

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username"`
	Password string        `bson:"password"`
	Fullname string        `bson:"fullname"`
	Mobile   string        `bson:"mobile"`
	UserType string        `bson:"userType"`
}
