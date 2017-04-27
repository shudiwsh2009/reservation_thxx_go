package web

import (
	"github.com/mijia/sweb/form"
	"github.com/mijia/sweb/render"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	"github.com/shudiwsh2009/reservation_thxx_go/service"
	"golang.org/x/net/context"
	"net/http"
)

type UserController struct {
	BaseMuxController
}

const (
	kUserApiBaseUrl = "/api/user"
)

func (uc *UserController) MuxHandlers(m JsonMuxer) {
	m.Get("/appointment", "EntryPageLegacy", uc.GetEntryPageLegacy)
	m.Get("/appointment/entry", "EntryPageLegacyBak", uc.GetEntryPageLegacy)
	m.Get("/appointment/student", "StudentPageLegacy", uc.GetStudentPageLegacy)
	m.Get("/appointment/teacher/login", "TeacherLoginPageLegacy", uc.GetTeacherLoginPageLegacy)
	m.Get("/appointment/teacher", "TeacherPageLegacy", LegacyPageInjection(uc.GetTeacherPageLegacy, "/appointment/teacher/login"))
	m.Get("/appointment/admin/login", "AdminLoginPageLegacy", uc.GetAdminLoginPageLegacy)
	m.Get("/appointment/admin", "AdminPageLegacy", LegacyPageInjection(uc.GetAdminPageLegacy, "/appointment/admin/login"))

	m.PostJson(kUserApiBaseUrl+"/teacher/login", "TeacherLogin", uc.TeacherLogin)
	m.PostJson(kUserApiBaseUrl+"/admin/login", "AdminLogin", uc.AdminLogin)
	m.PostJson(kUserApiBaseUrl+"/logout", "Logout", RoleCookieInjection(uc.Logout))
	m.PostJson(kUserApiBaseUrl+"/session", "UpdateSession", RoleCookieInjection(uc.UpdateSession))
}

func (uc *UserController) GetTemplates() []*render.TemplateSet {
	return []*render.TemplateSet{
		render.NewTemplateSet("entry", "desktop.html", "legacy/entry.html", "layout/desktop.html"),
		render.NewTemplateSet("student", "desktop.html", "legacy/student.html", "layout/desktop.html"),
		render.NewTemplateSet("teacher_login", "desktop.html", "legacy/teacher_login.html", "layout/desktop.html"),
		render.NewTemplateSet("teacher", "desktop.html", "legacy/teacher.html", "layout/desktop.html"),
		render.NewTemplateSet("admin_login", "desktop.html", "legacy/admin_login.html", "layout/desktop.html"),
		render.NewTemplateSet("admin", "desktop.html", "legacy/admin.html", "layout/desktop.html"),
	}
}

func (uc *UserController) GetEntryPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "entry", params)
	return ctx
}

func (uc *UserController) GetStudentPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "student", params)
	return ctx
}

func (uc *UserController) GetTeacherLoginPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "teacher_login", params)
	return ctx
}

func (uc *UserController) GetTeacherPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request, userId string, userType int) context.Context {
	if userType != model.UserTypeTeacher {
		http.Redirect(w, r, "/appointment/teacher/login", http.StatusFound)
		return ctx
	} else if teacher, err := service.MongoClient().GetTeacherById(userId); err != nil || teacher == nil || teacher.UserType != model.UserTypeTeacher {
		http.Redirect(w, r, "/appointment/teacher/login", http.StatusFound)
		return ctx
	}
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "teacher", params)
	return ctx
}

func (uc *UserController) GetAdminLoginPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "admin_login", params)
	return ctx
}

func (uc *UserController) GetAdminPageLegacy(ctx context.Context, w http.ResponseWriter, r *http.Request, userId string, userType int) context.Context {
	if userType != model.UserTypeAdmin {
		http.Redirect(w, r, "/appointment/admin/login", http.StatusFound)
		return ctx
	} else if admin, err := service.MongoClient().GetAdminById(userId); err != nil || admin == nil || admin.UserType != model.UserTypeAdmin {
		http.Redirect(w, r, "/appointment/admin/login", http.StatusFound)
		return ctx
	}
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "admin", params)
	return ctx
}

func (uc *UserController) TeacherLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().TeacherLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setSession(w, teacher.Id.Hex(), teacher.Username, teacher.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["redirect_url"] = "/appointment/teacher"

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) AdminLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := form.ParamString(r, "username", "")
	password := form.ParamString(r, "password", "")

	var result = make(map[string]interface{})

	admin, err := service.Workflow().AdminLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	if err = setSession(w, admin.Id.Hex(), admin.Username, admin.UserType); err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	result["redirect_url"] = "/appointment/admin"

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	switch userType {
	case model.UserTypeAdmin:
		result["redirect_url"] = "/appointment/admin/login"
	case model.UserTypeTeacher:
		result["redirect_url"] = "/appointment/teacher/login"
	default:
		result["redirect_url"] = "/appointment"
	}
	clearSession(w, r)

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) UpdateSession(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	result, err := service.Workflow().UpdateSession(userId, userType)
	if err != nil {
		return http.StatusOK, wrapJsonError(err)
	}
	return http.StatusOK, wrapJsonOk(result)
}
