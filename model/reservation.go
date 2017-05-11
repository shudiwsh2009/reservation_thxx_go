package model

import (
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	ReservationStatusAvailable  = 1
	ReservationStatusReservated = 2
	ReservationStatusFeedback   = 3
	ReservationStatusDeleted    = 4

	MakeReservationLastHour = 3 // 最迟提前3小时预约
)

type Reservation struct {
	Id              bson.ObjectId   `bson:"_id"`
	StartTime       time.Time       `bson:"start_time"`
	EndTime         time.Time       `bson:"end_time"`
	Status          int             `bson:"status"`
	TeacherUsername string          `bson:"teacher_username"`
	TeacherFullname string          `bson:"teacher_fullname"`
	TeacherMobile   string          `bson:"teacher_mobile"`
	TeacherAddress  string          `bson:"teacher_address"`
	StudentInfo     StudentInfo     `bson:"student_info"`
	StudentFeedback StudentFeedback `bson:"student_feedback"`
	TeacherFeedback TeacherFeedback `bson:"teacher_feedback"`
	CreatedAt       time.Time       `bson:"created_at"`
	UpdatedAt       time.Time       `bson:"updated_at"`
}

type StudentInfo struct {
	Fullname   string `bson:"fullname"`
	Gender     string `bson:"gender"`
	Username   string `bson:"username"`
	School     string `bson:"school`
	Hometown   string `bson:"hometown"`
	Mobile     string `bson:"mobile"`
	Email      string `bson:"email"`
	Experience string `bson:"experience"`
	Problem    string `bson:"problem"`
}

type StudentFeedback struct {
	Fullname string `bson:"fullname"`
	Problem  string `bson:"problem"`
	Choices  string `bson:"choices"`
	Score    string `bson:"score"`
	Feedback string `bson:"feedback"`
}

func (sf StudentFeedback) IsEmpty() bool {
	return sf.Fullname == "" || sf.Problem == "" || sf.Choices == "" || sf.Score == "" || sf.Feedback == ""
}

type TeacherFeedback struct {
	TeacherFullname string `bson:"teacher_fullname"`
	TeacherUsername string `bson:"teacher_username"`
	StudentFullname string `bson:"student_fullname"`
	Problem         string `bson:"problem"`
	Solution        string `bson:"solution"`
	AdviceToCenter  string `bson:"advice_to_center"`
}

func (tf TeacherFeedback) IsEmpty() bool {
	return tf.TeacherFullname == "" || tf.TeacherUsername == "" || tf.StudentFullname == "" ||
		tf.Problem == "" || tf.Solution == "" || tf.AdviceToCenter == ""
}

func (m *MongoClient) InsertReservation(reservation *Reservation) error {
	now := time.Now()
	reservation.Id = bson.NewObjectId()
	reservation.CreatedAt = now
	reservation.UpdatedAt = now
	return dbReservation.Insert(reservation)
}

func (m *MongoClient) UpdateReservation(reservation *Reservation) error {
	reservation.UpdatedAt = time.Now()
	return dbReservation.UpdateId(reservation.Id, reservation)
}

func (m *MongoClient) UpdateReservationWithoutUpdatedTime(reservation *Reservation) error {
	return dbReservation.UpdateId(reservation.Id, reservation)
}

func (m *MongoClient) GetAllReservations() ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{}).All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservationById(id string) (*Reservation, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ErrorDatabase)
	}
	var reservation Reservation
	err := dbReservation.FindId(bson.ObjectIdHex(id)).One(&reservation)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &reservation, nil
	}
}

func (m *MongoClient) GetReservationsByStudentUsername(studentUsername string) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"student_info.username": studentUsername,
		"status": bson.M{"$ne": ReservationStatusDeleted}}).Sort("start_time").All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservationsBetweenTime(start time.Time, end time.Time) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"start_time": bson.M{"$gte": start, "$lt": end},
		"status": bson.M{"$ne": ReservationStatusDeleted}}).Sort("start_time").All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservatedReservationsBetweenTime(start time.Time, end time.Time) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"start_time": bson.M{"$gte": start, "$lt": end},
		"status": ReservationStatusReservated}).Sort("start_time").All(&reservations)
	return reservations, err
}

func (m *MongoClient) GetReservationsAfterTime(start time.Time) ([]*Reservation, error) {
	var reservations []*Reservation
	err := dbReservation.Find(bson.M{"start_time": bson.M{"$gte": start},
		"status": bson.M{"$ne": ReservationStatusDeleted}}).Sort("start_time").All(&reservations)
	return reservations, err
}
