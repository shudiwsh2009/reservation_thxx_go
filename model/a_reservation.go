package model

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ReservationStatus int

const (
	AVAILABLE ReservationStatus = 1 + iota
	RESERVATED
	FEEDBACK
	DELETED
)

var reservationStatuses = [...]string{
	"AVAILABLE",
	"RESERVATED",
	"FEEDBACK",
	"DELETED",
}

func (rs ReservationStatus) String() string {
	return reservationStatuses[rs-1]
}

type LegacyStudentInfo struct {
	Name       string `bson:"name"`
	Gender     string `bson:"gender"`
	StudentId  string `bson:"studentId"`
	School     string `bson:"school`
	Hometown   string `bson:"hometown"`
	Mobile     string `bson:"mobile"`
	Email      string `bson:"email"`
	Experience string `bson:"experience"`
	Problem    string `bson:"problem"`
}

type LegacyStudentFeedback struct {
	Name     string `bson:"name"`
	Problem  string `bson:"problem"`
	Choices  string `bson:"choices"`
	Score    string `bson:"score"`
	Feedback string `bson:"feedback"`
}

func (sf LegacyStudentFeedback) IsEmpty() bool {
	return len(sf.Name) == 0 || len(sf.Problem) == 0 || len(sf.Choices) == 0 ||
		len(sf.Score) == 0 || len(sf.Feedback) == 0
}

type LegacyTeacherFeedback struct {
	TeacherFullname string `bson:"teacherName"`
	TeacherUsername string `bson:"teacherId"`
	StudentFullname string `bson:"studentName"`
	Problem         string `bson:"problem"`
	Solution        string `bson:"solution"`
	AdviceToCenter  string `bson:"adviceToCenter"`
}

func (tf LegacyTeacherFeedback) IsEmpty() bool {
	return len(tf.TeacherFullname) == 0 || len(tf.TeacherUsername) == 0 || len(tf.StudentFullname) == 0 ||
		len(tf.Problem) == 0 || len(tf.Solution) == 0 || len(tf.AdviceToCenter) == 0
}

type LegacyReservation struct {
	Id              bson.ObjectId         `bson:"_id"`
	StartTime       time.Time             `bson:"startTime"`
	EndTime         time.Time             `bson:"endTime"`
	Status          ReservationStatus     `bson:"status_GO"`
	TeacherUsername string                `bson:"teacherUsername"`
	TeacherFullname string                `bson:"teacher"`
	TeacherMobile   string                `bson:"teacherMobile"`
	TeacherAddress  string                `bson:"teacherAddress"`
	StudentInfo     LegacyStudentInfo     `bson:"studentInfo"`
	StudentFeedback LegacyStudentFeedback `bson:"studentFeedback"`
	TeacherFeedback LegacyTeacherFeedback `bson:"teacherFeedback"`
}

func (m *LegacyMongoClient) GetAllLegacyReservations() ([]*LegacyReservation, error) {
	var reservations []*LegacyReservation
	err := dbLegacyReservation.Find(bson.M{}).All(&reservations)
	return reservations, err
}
