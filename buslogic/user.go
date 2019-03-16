package buslogic

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"time"
)

const (
	AdminDefaultPassword   = "THXXFZZX"
	TeacherDefaultPassword = "thxxfzzx"
)

var (
	UserLoginMaxCount = map[int]int{
		model.UserTypeAdmin:   1,
		model.UserTypeTeacher: 1,
		model.UserTypeStudent: 5,
	}
	UserLoginExpire = map[int]time.Duration{
		model.UserTypeAdmin:   time.Hour * 24 * 7,
		model.UserTypeTeacher: time.Hour * 24 * 7,
		model.UserTypeStudent: time.Hour * 24 * 30,
	}
)

// 咨询师登录
func (w *Workflow) TeacherLogin(username string, password string) (*model.Teacher, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ErrorMissingParam, "password")
	}
	teacher, err := w.MongoClient().GetTeacherByUsername(username)
	if err == nil && teacher != nil && teacher.Password == model.EncodePassword(teacher.Salt, password) {
		return teacher, nil
	}
	return nil, re.NewRErrorCode("wrong password", nil, re.ErrorLoginPasswordWrong)
}

// 管理员登录
func (w *Workflow) AdminLogin(username string, password string) (*model.Admin, error) {
	if username == "" {
		return nil, re.NewRErrorCodeContext("username is empty", nil, re.ErrorMissingParam, "username")
	} else if password == "" {
		return nil, re.NewRErrorCodeContext("password is empty", nil, re.ErrorMissingParam, "password")
	}
	admin, err := w.MongoClient().GetAdminByUsername(username)
	if err == nil && admin != nil && admin.Password == model.EncodePassword(admin.Salt, password) {
		return admin, nil
	}
	return nil, re.NewRErrorCode("wrong password", nil, re.ErrorLoginPasswordWrong)
}

// 更新session
func (w *Workflow) UpdateSession(userId string, userType int) (map[string]interface{}, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("user not login", nil, re.ErrorNoLogin)
	}
	result := make(map[string]interface{})
	switch userType {
	case model.UserTypeAdmin:
		admin, err := w.MongoClient().GetAdminById(userId)
		if err != nil || admin == nil || admin.UserType != userType {
			return nil, re.NewRErrorCode("fail to get admin", err, re.ErrorDatabase)
		}
		//result["user"] = w.WrapAdmin(admin)
	case model.UserTypeTeacher:
		teacher, err := w.MongoClient().GetTeacherById(userId)
		if err != nil || teacher == nil || teacher.UserType != userType {
			return nil, re.NewRErrorCode("fail to get teacher", err, re.ErrorDatabase)
		}
		result["teacher"] = w.WrapTeacher(teacher)
	default:
		return nil, re.NewRErrorCode("fail to get user", nil, re.ErrorNoUser)
	}
	return result, nil
}

// external: 重置账户密码
func (w *Workflow) ResetUserPassword(username string, userType int, password string) error {
	if username == "" || password == "" {
		return re.NewRError("missing parameters", nil)
	}
	var err error
	var userId string
	switch userType {
	case model.UserTypeTeacher:
		teacher, err := w.MongoClient().GetTeacherByUsername(username)
		if err != nil || teacher == nil || teacher.UserType != userType {
			return re.NewRError("fail to get teacher", err)
		}
		teacher.Password = password
		teacher.PreInsert()
		err = w.MongoClient().UpdateTeacher(teacher)
		userId = teacher.Id.Hex()
	case model.UserTypeAdmin:
		admin, err := w.MongoClient().GetAdminByUsername(username)
		if err != nil || admin == nil || admin.UserType != userType {
			return re.NewRError("fail to get admin", err)
		}
		admin.Password = password
		admin.PreInsert()
		err = w.MongoClient().UpdateAdmin(admin)
		userId = admin.Id.Hex()
	default:
		return re.NewRError(fmt.Sprintf("unknown user_type: %d", userType), nil)
	}
	if err != nil {
		return re.NewRError("fail to update user", err)
	}
	return w.ClearUserLoginRedisKey(userId, userType)
}

func (w *Workflow) ClearUserLoginRedisKey(userId string, userType int) error {
	redisKeys, err := w.RedisClient().Keys(fmt.Sprintf(model.RedisKeyLogin, userType, userId, "*")).Result()
	if err != nil {
		return re.NewRError("fail to get user login session keys from redis", err)
	}
	for _, k := range redisKeys {
		if err := w.RedisClient().Del(k).Err(); err != nil {
			return err
		}
	}
	return nil
}

// external: 添加新管理员
func (w *Workflow) AddNewAdmin(username string, password string) (*model.Admin, error) {
	if username == "" || password == "" {
		return nil, re.NewRError("missing parameters", nil)
	}
	oldAdmin, err := w.MongoClient().GetAdminByUsername(username)
	if err != nil {
		return nil, re.NewRError("fail to get old admin", err)
	} else if oldAdmin != nil && oldAdmin.UserType == model.UserTypeAdmin {
		return oldAdmin, re.NewRError(fmt.Sprintf("admin already exists: %+v", oldAdmin), nil)
	}
	newAdmin := &model.Admin{
		Username: username,
		Password: password,
		UserType: model.UserTypeAdmin,
	}
	err = w.MongoClient().InsertAdmin(newAdmin)
	return newAdmin, err
}
