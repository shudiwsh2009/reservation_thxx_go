package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/data"
	"github.com/shudiwsh2009/reservation_thxx_go/domain"
	"strings"
)

const (
	AdminDefaultPassword   = "THXXFZZX"
	TeacherDefaultPassword = "thxxfzzx"
)

type UserLogic struct {
}

// 学生登录
func (ul *UserLogic) Login(username string, password string) (*domain.User, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("用户名为空")
	} else if strings.EqualFold(password, "") {
		return nil, errors.New("密码为空")
	}
	user, err := data.GetUserByUsername(username)
	if err == nil && (strings.EqualFold(user.Password, password) ||
		(strings.EqualFold(user.UserType, domain.TEACHER) && strings.EqualFold(user.Password, TeacherDefaultPassword)) ||
		(strings.EqualFold(user.UserType, domain.ADMIN) && strings.EqualFold(user.Password, AdminDefaultPassword))) {
		return user, nil
	}
	return nil, errors.New("用户名或密码不正确")
}

// 获取用户
func (ul *UserLogic) GetUserByUsername(username string) (*domain.User, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	}
	user, err := data.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// 查找咨询师
// 查找顺序:全名 > 工号 > 手机号
func (ul *UserLogic) SearchTeacher(fullname string, username string, mobile string) (*domain.User, error) {
	if !strings.EqualFold(fullname, "") {
		user, err := data.GetUserByFullname(fullname)
		if err == nil {
			return user, nil
		}
	}
	if !strings.EqualFold(username, "") {
		user, err := data.GetUserByUsername(username)
		if err == nil {
			return user, nil
		}
	}
	if !strings.EqualFold(mobile, "") {
		user, err := data.GetUserByMobile(mobile)
		if err == nil {
			return user, nil
		}
	}
	return nil, errors.New("用户不存在")
}
