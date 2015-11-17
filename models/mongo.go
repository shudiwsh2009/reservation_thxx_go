package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	Mongo *mgo.Database
)

/**
User
*/
func AddSimpleUser(username string, password string, userType UserType) (*User, error) {
	collection := Mongo.C("user")
	newUser := &User{
		Id:       bson.NewObjectId(),
		Username: username,
		Password: password,
		UserType: userType,
	}
	if err := collection.Insert(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func AddFullUser(username string, password string, fullname string, mobile string, userType UserType) (*User, error) {
	collection := Mongo.C("user")
	newUser := &User{
		Id:       bson.NewObjectId(),
		Username: username,
		Password: password,
		Fullname: fullname,
		Mobile:   mobile,
		UserType: userType,
	}
	if err := collection.Insert(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func UpsertUser(user *User) error {
	collection := Mongo.C("user")
	_, err := collection.UpsertId(user.Id, user)
	return err
}

func GetUserById(userId string) (*User, error) {
	collection := Mongo.C("user")
	user := &User{}
	if err := collection.FindId(bson.ObjectIdHex(userId)).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	collection := Mongo.C("user")
	user := &User{}
	if err := collection.Find(bson.M{"username": username}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByFullname(fullname string) (*User, error) {
	collection := Mongo.C("user")
	user := &User{}
	if err := collection.Find(bson.M{"fullname": fullname}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByMobile(mobile string) (*User, error) {
	collection := Mongo.C("user")
	user := &User{}
	if err := collection.Find(bson.M{"mobile": mobile}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

/**
Reservation
*/

func AddReservation(startTime time.Time, endTime time.Time, teacherFullname string, teacherUsername string,
	teacherMobile string) (*Reservation, error) {
	collection := Mongo.C("appointment")
	newReservation := &Reservation{
		Id:              bson.NewObjectId(),
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          AVAILABLE,
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherUsername,
		TeacherMobile:   teacherMobile,
		StudentInfo:     StudentInfo{},
		StudentFeedback: StudentFeedback{},
		TeacherFeedback: TeacherFeedback{},
	}
	if err := collection.Insert(newReservation); err != nil {
		return nil, err
	}
	return newReservation, nil
}

func UpsertReservation(reservation *Reservation) error {
	collection := Mongo.C("appointment")
	_, err := collection.UpsertId(reservation.Id, reservation)
	return err
}

func GetReservationById(id string) (*Reservation, error) {
	collection := Mongo.C("appointment")
	reservation := &Reservation{}
	if err := collection.FindId(bson.ObjectIdHex(id)).One(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func GetReservationsByStudentId(studentId string) ([]*Reservation, error) {
	collection := Mongo.C("appointment")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"studentInfo.studentId": studentId}).All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func GetReservationsBetweenTime(from time.Time, to time.Time) ([]*Reservation, error) {
	collection := Mongo.C("appointment")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"startTime": bson.M{"$gte": from, "$lte": to},
		"status": bson.M{"$ne": DELETED}}).Sort("startTime").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func GetReservationsAfterTime(from time.Time) ([]*Reservation, error) {
	collection := Mongo.C("appointment")
	var reservations []*Reservation
	if err := collection.Find(bson.M{"startTime": bson.M{"$gte": from},
		"status": bson.M{"$ne": DELETED}}).Sort("startTime").All(&reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}
