package model

import (
	"crypto/sha256"
	"encoding/base64"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	UserTypeUnknown = 0
	UserTypeStudent = 1
	UserTypeTeacher = 2
	UserTypeAdmin   = 3

	InternationalTypeChinese   = 0
	InternationalTypeChinglish = 1

	GraduateTypeBoth  = 0 // 均可
	GraduateTypeUnder = 1 // 仅限本科生
	GraduateTypePost  = 2 // 仅限研究生
)

type Teacher struct {
	Id                bson.ObjectId `bson:"_id"`
	Username          string        `bson:"username"` // Indexed
	Password          string        `bson:"password"`
	Salt              string        `bson:"salt"`
	UserType          int           `bson:"user_type"`
	Fullname          string        `bson:"fullname"`
	FullnameEn        string        `bson:"fullname_en"`
	Mobile            string        `bson:"mobile"`
	Address           string        `bson:"address"`
	AddressEn         string        `bson:"address_en"`
	Gender            string        `bson:"gender"`
	GenderEn          string        `bson:"gender_en"`
	Major             string        `bson:"major"` // 专业背景
	MajorEn           string        `bson:"major_en"`
	Academic          string        `bson:"academic"` // 学历
	AcademicEn        string        `bson:"academic_en"`
	Aptitude          string        `bson:"aptitude"` // 资质
	AptitudeEn        string        `bson:"aptitude_en"`
	Problem           string        `bson:"problem"` // 可咨询的问题
	ProblemEn         string        `bson:"problem_en"`
	SmsSuffix         string        `bson:"sms_suffix"` // 发送短信自动添加后缀
	SmsSuffixEn       string        `bson:"sms_suffix_en"`
	InternationalType int           `bson:"international_type"` // 国际化类型：0、仅中文 1、中英双语
	GraduateType      int           `bson:"graduate_type"` // 学生类型：0、均可 1、仅本科生 2、仅研究生
	CreatedAt         time.Time     `bson:"created_at"`
	UpdatedAt         time.Time     `bson:"updated_at"`
}

func (teacher *Teacher) PreInsert() error {
	salt := EncodePassword("salt", strconv.Itoa(rand.Int()))
	teacher.Salt = salt[:16]
	teacher.Password = EncodePassword(teacher.Salt, teacher.Password)
	teacher.Username = strings.TrimSpace(teacher.Username)
	teacher.UserType = UserTypeTeacher
	return nil
}

func (m *MongoClient) InsertTeacher(teacher *Teacher) error {
	teacher.PreInsert()
	now := time.Now()
	teacher.Id = bson.NewObjectId()
	teacher.CreatedAt = now
	teacher.UpdatedAt = now
	return dbTeacher.Insert(teacher)
}

func (m *MongoClient) UpdateTeacher(teacher *Teacher) error {
	teacher.UpdatedAt = time.Now()
	return dbTeacher.UpdateId(teacher.Id, teacher)
}

func (m *MongoClient) UpdateTeacherWithoutTime(teacher *Teacher) error {
	return dbTeacher.UpdateId(teacher.Id, teacher)
}

func (m *MongoClient) GetTeacherById(id string) (*Teacher, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ErrorDatabase)
	}
	var teacher Teacher
	err := dbTeacher.FindId(bson.ObjectIdHex(id)).One(&teacher)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &teacher, nil
	}
}

func (m *MongoClient) GetTeacherByUsername(username string) (*Teacher, error) {
	var teacher Teacher
	err := dbTeacher.Find(bson.M{"username": username, "user_type": UserTypeTeacher}).One(&teacher)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &teacher, nil
	}
}

func (m *MongoClient) GetTeacherByFullname(fullname string) (*Teacher, error) {
	var teacher Teacher
	err := dbTeacher.Find(bson.M{"fullname": fullname, "user_type": UserTypeTeacher}).One(&teacher)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &teacher, nil
	}
}

func (m *MongoClient) GetTeacherByFullnameEn(fullnameEn string) (*Teacher, error) {
	var teacher Teacher
	err := dbTeacher.Find(bson.M{"fullname_en": fullnameEn, "user_type": UserTypeTeacher}).One(&teacher)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &teacher, nil
	}
}

func (m *MongoClient) GetTeacherByMobile(mobile string) (*Teacher, error) {
	var teacher Teacher
	err := dbTeacher.Find(bson.M{"mobile": mobile, "user_type": UserTypeTeacher}).One(&teacher)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &teacher, nil
	}
}

type Admin struct {
	Id        bson.ObjectId `bson:"_id"`
	Username  string        `bson:"username"`
	Password  string        `bson:"password"`
	Salt      string        `bson:"salt"`
	UserType  int           `bson:"user_type"`
	Fullname  string        `bson:"fullname"`
	Mobile    string        `bson:"mobile"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (admin *Admin) PreInsert() error {
	salt := EncodePassword("salt", strconv.Itoa(rand.Int()))
	admin.Salt = salt[:16]
	admin.Password = EncodePassword(admin.Salt, admin.Password)
	admin.Username = strings.TrimSpace(admin.Username)
	admin.UserType = UserTypeAdmin
	return nil
}

func (m *MongoClient) InsertAdmin(admin *Admin) error {
	admin.PreInsert()
	now := time.Now()
	admin.Id = bson.NewObjectId()
	admin.CreatedAt = now
	admin.UpdatedAt = now
	return dbAdmin.Insert(admin)
}

func (m *MongoClient) UpdateAdmin(admin *Admin) error {
	admin.UpdatedAt = time.Now()
	return dbAdmin.UpdateId(admin.Id, admin)
}

func (m *MongoClient) UpdateAdminWithoutTime(admin *Admin) error {
	return dbAdmin.UpdateId(admin.Id, admin)
}

func (m *MongoClient) GetAdminById(id string) (*Admin, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, re.NewRErrorCode("id is not valid", nil, re.ErrorDatabase)
	}
	var admin Admin
	err := dbAdmin.FindId(bson.ObjectIdHex(id)).One(&admin)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &admin, nil
	}
}

func (m *MongoClient) GetAdminByUsername(username string) (*Admin, error) {
	var admin Admin
	err := dbAdmin.Find(bson.M{"username": username, "user_type": UserTypeAdmin}).One(&admin)
	if err == mgo.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &admin, nil
	}
}

// 对密码加盐加密
func EncodePassword(salt, passwd string) string {
	h := sha256.New()
	h.Write([]byte(passwd + salt))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
