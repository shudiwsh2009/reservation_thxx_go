package data

import (
	"github.com/shudiwsh2009/reservation_thxx_go/domain"
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
func AddSimpleUser(username string, password string, userType domain.UserType) (*domain.User, error) {
	collection := Mongo.C("user")
	newUser := &domain.User{
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

func AddFullUser(username string, password string, fullname string, mobile string, userType domain.UserType) (*domain.User, error) {
	collection := Mongo.C("user")
	newUser := &domain.User{
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

func UpsertUser(user *domain.User) error {
	collection := Mongo.C("user")
	_, err := collection.UpsertId(user.Id, user)
	return err
}

func GetUserByUsername(username string) (*domain.User, error) {
	collection := Mongo.C("user")
	user := &domain.User{}
	if err := collection.Find(bson.M{"username": username}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByFullname(fullname string) (*domain.User, error) {
	collection := Mongo.C("user")
	user := &domain.User{}
	if err := collection.Find(bson.M{"fullname": fullname}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByMobile(mobile string) (*domain.User, error) {
	collection := Mongo.C("user")
	user := &domain.User{}
	if err := collection.Find(bson.M{"mobile": mobile}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

/**
Reservation
*/

func AddReservation(startTime time.Time, endTime time.Time, teacherFullname string, teacherUsername string,
	teacherMobile string) (*domain.Reservation, error) {
	collection := Mongo.C("appointment")
	newReservation := &domain.Reservation{
		Id:              bson.NewObjectId(),
		StartTime:       startTime,
		EndTime:         endTime,
		Status:          domain.Availabel,
		TeacherFullname: teacherFullname,
		TeacherUsername: teacherUsername,
		TeacherMobile:   teacherMobile,
		StudentInfo:     domain.StudentInfo{},
		StudentFeedback: domain.StudentFeedback{},
		TeacherFeedback: domain.TeacherFeedback{},
	}
	if err := collection.Insert(newReservation); err != nil {
		return nil, err
	}
	return newReservation, nil
}

func UpsertReservation(reservation *domain.Reservation) error {
	collection := Mongo.C("appointment")
	_, err := collection.UpsertId(reservation.Id, reservation)
	return err
}

func GetReservationById(id string) (*domain.Reservation, error) {
	collection := Mongo.C("appointment")
	reservation := &domain.Reservation{}
	if err := collection.FindId(bson.ObjectIdHex(id)).One(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func GetReservationsByStudentId(studentId string) ([]*domain.Reservation, error) {
	collection := Mongo.C("appointment")
	var reservations []*domain.Reservation
	reservation := &domain.Reservation{}
	iter := collection.Find(bson.M{"studentInfo.studentId": studentId}).Iter()
	for iter.Next(reservation) {
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func GetReservationsBetweenTime(from time.Time, to time.Time) ([]*domain.Reservation, error) {
	collection := Mongo.C("appointment")
	var reservations []*domain.Reservation
	reservation := &domain.Reservation{}
	iter := collection.Find(bson.M{"startTime": bson.M{"$gte": from, "$lte": to},
		"status": bson.M{"$ne": domain.Deleted}}).Sort("startTime").Iter()
	for iter.Next(reservation) {
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}

func GetReservationsAfterTime(from time.Time) ([]*domain.Reservation, error) {
	collection := Mongo.C("appointment")
	var reservations []*domain.Reservation
	reservation := &domain.Reservation{}
	iter := collection.Find(bson.M{"startTime": bson.M{"$gte": from},
		"status": bson.M{"$ne": domain.Deleted}}).Sort("startTime").Iter()
	for iter.Next(reservation) {
		reservations = append(reservations, reservation)
	}
	return reservations, nil
}
