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
	if len(username) == 0 {
		return nil, errors.New("用户名为空")
	} else if len(password) == 0 {
		return nil, errors.New("密码为空")
	}
	user, err := models.GetUserByUsername(username)
	if err == nil && (strings.EqualFold(password, user.Password) ||
		(user.UserType == models.TEACHER && strings.EqualFold(password, TeacherDefaultPassword)) ||
		(user.UserType == models.ADMIN && strings.EqualFold(password, AdminDefaultPassword))) {
		return user, nil
	}
	return nil, errors.New("用户名或密码不正确")
}

// 获取用户
func (ul *UserLogic) GetUserByUsername(username string) (*models.User, error) {
	if len(username) == 0 {
		return nil, errors.New("请先登录")
	}
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

func (ul *UserLogic) GetUserById(userId string) (*models.User, error) {
	if len(userId) == 0 {
		return nil, errors.New("请先登录")
	}
	user, err := models.GetUserById(userId)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}
