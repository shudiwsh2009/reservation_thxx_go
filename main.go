package main

import (
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/data"
	"gopkg.in/mgo.v2"
	"time"
)

func main() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		fmt.Errorf("连接数据库失败:%v", err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	data.Mongo = session.DB("appointment")

	t, err := time.Parse(time.UnixDate, "Sat Mar  7 14:06:39 PST 2015")
	if err != nil { // Always check errors even if they should not happen.
		panic(err)
	}
	do := func(name, layout, want string) {
		got := t.Format(layout)
		if want != got {
			fmt.Printf("error: for %q got %q; expected %q\n", layout, got, want)
			return
		}
		fmt.Printf("%-15s %q gives %q\n", name, layout, got)
	}
	do("Suppressed pad", "15:04", "14:06")
	do("month", "一", "三")
}
