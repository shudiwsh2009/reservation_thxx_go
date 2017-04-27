package web

import (
	"github.com/mijia/sweb/form"
	"github.com/mijia/sweb/server"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	"github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/service"
	"golang.org/x/net/context"
	"net/http"
)

func RoleCookieInjection(handle func(http.ResponseWriter, *http.Request, string, int) (int, interface{})) JsonHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
		userId, _, userType, err := getSession(r)
		if err != nil {
			return http.StatusOK, wrapJsonError(err)
		}
		return handle(w, r, userId, userType)
	}
}

func FakeInjection(handle func(http.ResponseWriter, *http.Request, string, int) (int, interface{})) JsonHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
		//userId, _, userType, err := getSession(r)
		//if err != nil {
		//	return http.StatusOK, wrapJsonError(err)
		//}
		return handle(w, r, "", model.UserTypeUnknown)
	}
}

func LegacyPageInjection(handle func(context.Context, http.ResponseWriter, *http.Request, string, int) context.Context, redirectUrl string ) server.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
		userId, _, userType, err := getSession(r)
		if err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
		}
		return handle(ctx, w, r, userId, userType)
	}
}

func RequestPasswordCheck(r *http.Request, userId string, userType int) error {
	password := form.ParamString(r, "password", "")
	if password == "" {
		return rerror.NewRErrorCodeContext("password is empty", nil, rerror.ErrorMissingParam, "password")
	}
	switch userType {
	case model.UserTypeTeacher:
		teacher, err := service.MongoClient().GetTeacherById(userId)
		if err == nil && teacher != nil && teacher.UserType == model.UserTypeTeacher {
			if teacher.Password == model.EncodePassword(teacher.Salt, password) {
				return nil
			}
		}
	case model.UserTypeAdmin:
		admin, err := service.MongoClient().GetAdminById(userId)
		if err == nil && admin != nil && admin.UserType == model.UserTypeAdmin {
			if admin.Salt != "" && admin.Password == model.EncodePassword(admin.Salt, password) {
				return nil
			}
		}
	}
	return rerror.NewRErrorCode("request password check failed", nil, rerror.ErrorNotAuthorized)
}
