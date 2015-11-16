package buslogic

import (
	"errors"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"strings"
)

const (
	AdminDefaultPassword   = "THXXFZZX"
	TeacherDefaultPassword = "thxxfzzx"
)

type UserLogic struct {
}

// 学生登录
func (ul *UserLogic) Login(username string, password string) (*models.User, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("用户名为空")
	} else if strings.EqualFold(password, "") {
		return nil, errors.New("密码为空")
	}
	user, err := models.GetUserByUsername(username)
	if err == nil && (strings.EqualFold(user.Password, password) ||
		(user.UserType == models.TEACHER && strings.EqualFold(user.Password, TeacherDefaultPassword)) ||
		(user.UserType == models.ADMIN && strings.EqualFold(user.Password, AdminDefaultPassword))) {
		return user, nil
	}
	return nil, errors.New("用户名或密码不正确")
}

// 获取用户
func (ul *UserLogic) GetUserByUsername(username string) (*models.User, error) {
	if strings.EqualFold(username, "") {
		return nil, errors.New("请先登录")
	}
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}
