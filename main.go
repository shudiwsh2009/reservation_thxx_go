package main

import (
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/controllers"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

var needUserPath = regexp.MustCompile("^/reservation/(logout|(teacher|admin)/)")

func checkUserHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check url to see whether there is "/teacher/" or "/admin/" or "/logout"
		m := needUserPath.FindStringSubmatch(r.URL.Path)
		fmt.Println(r.URL.Path, "   ", m)
		if m == nil {
			fn(w, r)
			return
		}
		if _, err := r.Cookie("user_id"); err != nil {
			controllers.ErrorHandler(w, r, errors.New("请先登录"))
			return
		} else if _, err := r.Cookie("username"); err != nil {
			controllers.ErrorHandler(w, r, errors.New("请先登录"))
			return
		} else if _, err := r.Cookie("user_type"); err != nil {
			controllers.ErrorHandler(w, r, errors.New("请先登录"))
			return
		}
		fn(w, r)
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
	// 加载静态资源处理器
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	// 加载页面处理器
	http.HandleFunc("/appointment", controllers.EntryPage)
	http.HandleFunc("/appointment/entry", controllers.EntryPage)
	http.HandleFunc("/appointment/login", controllers.LoginPage)
	http.HandleFunc("/appointment/student", controllers.StudentPage)
	http.HandleFunc("/appointment/teacher", controllers.TeacherPage)
	http.HandleFunc("/appointment/admin", controllers.AdminPage)
	// 加载动态处理器
	http.HandleFunc("/reservation/login", checkUserHandler(controllers.Login))
	http.HandleFunc("/reservation/logout", checkUserHandler(controllers.Logout))
	http.HandleFunc("/reservation/teacher/view", checkUserHandler(controllers.ViewReservationsByStudent))
	// 启动监听
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
