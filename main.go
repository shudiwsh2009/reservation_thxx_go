package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shudiwsh2009/reservation_thxx_go/controllers"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var needUserPath = regexp.MustCompile("^(/appointment/(teacher|admin)|/(user/logout|(teacher|admin)/))")

func handleWithCookie(fn func(http.ResponseWriter, *http.Request, string, models.UserType)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check url to see whether there is "/teacher/" or "/admin/" or "/logout"
		m := needUserPath.FindStringSubmatch(r.URL.Path)
		if len(m) == 0 {
			fn(w, r, "", 0)
			return
		}
		var userId string
		var userType models.UserType
		if cookie, err := r.Cookie("user_id"); err != nil {
			http.Redirect(w, r, "/appointment/login", http.StatusFound)
			return
		} else {
			userId = cookie.Value
		}
		if _, err := r.Cookie("username"); err != nil {
			http.Redirect(w, r, "/appointment/login", http.StatusFound)
			return
		}
		if cookie, err := r.Cookie("user_type"); err != nil {
			http.Redirect(w, r, "/appointment/login", http.StatusFound)
			return
		} else {
			if ut, err := strconv.Atoi(cookie.Value); err != nil {
				http.Redirect(w, r, "/appointment/login", http.StatusFound)
				return
			} else {
				userType = models.UserType(ut)
			}
		}
		fn(w, r, userId, userType)
	}
}

func main() {
	// 数据库连接
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		fmt.Errorf("连接数据库失败：%v", err)
		os.Exit(1)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("appointment")
	// 时区
	if utils.Location, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		fmt.Errorf("初始化时区失败：%v", err)
		os.Exit(1)
	}

	// TODO: Remove the following test codes

	// mux
	router := mux.NewRouter()
	// 加载页面处理器
	pageRouter := router.PathPrefix("/appointment").Methods("GET").Subrouter()
	pageRouter.HandleFunc("", handleWithCookie(controllers.EntryPage))
	pageRouter.HandleFunc("/entry", handleWithCookie(controllers.EntryPage))
	pageRouter.HandleFunc("/login", handleWithCookie(controllers.LoginPage))
	pageRouter.HandleFunc("/student", handleWithCookie(controllers.StudentPage))
	pageRouter.HandleFunc("/teacher", handleWithCookie(controllers.TeacherPage))
	pageRouter.HandleFunc("/admin", handleWithCookie(controllers.AdminPage))
	// 加载动态处理器
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", handleWithCookie(controllers.Login)).Methods("POST")
	userRouter.HandleFunc("/logout", handleWithCookie(controllers.Logout)).Methods("GET")
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/reservation/view", handleWithCookie(controllers.ViewReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/view/monthly", handleWithCookie(controllers.ViewMonthlyReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/add", handleWithCookie(controllers.AddReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/edit", handleWithCookie(controllers.EditReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/remove", handleWithCookie(controllers.RemoveReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/cancel", handleWithCookie(controllers.CancelReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/feedback/get", handleWithCookie(controllers.GetFeedbackByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/feedback/submit", handleWithCookie(controllers.SubmitFeedbackByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/export", handleWithCookie(controllers.ExportReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/get", handleWithCookie(controllers.GetStudentInfoByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/teacher/search", handleWithCookie(controllers.SearchTeacherByAdmin)).Methods("POST")
	// http加载处理器
	http.Handle("/", router)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
