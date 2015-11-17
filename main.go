package main

import (
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
	"github.com/gorilla/mux"
)

var needUserPath = regexp.MustCompile("^(/appointment/(teacher|admin)|/reservation/(user/logout|(teacher|admin)/))")

func checkUser(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check url to see whether there is "/teacher/" or "/admin/" or "/logout"
		m := needUserPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fn(w, r)
			return
		}
		fmt.Println("check user cookie in URL: ", r.URL.Path)
		if _, err := r.Cookie("user_id"); err != nil {
			http.Redirect(w, r, "/appointment/login", http.StatusFound)
			return
		} else if _, err := r.Cookie("username"); err != nil {
			http.Redirect(w, r, "/appointment/login", http.StatusFound)
			return
		} else if _, err := r.Cookie("user_type"); err != nil {
			http.Redirect(w, r, "/appointment/login", http.StatusFound)
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

	// TODO: Remove the following test codes

	// mux
	router := mux.NewRouter()
	// 加载页面处理器
	pageRouter := router.PathPrefix("/appointment").Methods("GET").Subrouter()
	pageRouter.HandleFunc("", checkUser(controllers.EntryPage))
	pageRouter.HandleFunc("/entry", checkUser(controllers.EntryPage))
	pageRouter.HandleFunc("/login", checkUser(controllers.LoginPage))
	pageRouter.HandleFunc("/student", checkUser(controllers.StudentPage))
	pageRouter.HandleFunc("/teacher", checkUser(controllers.TeacherPage))
	pageRouter.HandleFunc("/admin", checkUser(controllers.AdminPage))
	// 加载动态处理器
	dynamicRouter := router.PathPrefix("/reservation").Methods("POST").Subrouter()
	userRouter := dynamicRouter.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/login", checkUser(controllers.Login))
	userRouter.HandleFunc("/logout", checkUser(controllers.Logout))
	// http加载处理器
	http.Handle("/", router)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
