package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/buslogic"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/http"
)

func ViewReservationsByStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var result = map[string]interface{}{"state": "SUCCESS"}
		var rl = buslogic.ReservationLogic{}

		reservations, err := rl.GetReservationsByStudent()
		if err != nil {
			ErrorHandler(w, r, err)
			return
		}
		var array []interface{}
		for _, res := range reservations {
			object := make(map[string]interface{})
			object["reservation_id"] = res.Id
			object["start_time"] = res.StartTime.Format(utils.TIME_PATTERN)
			object["end_time"] = res.EndTime.Format(utils.TIME_PATTERN)
			object["teacher_fullname"] = res.TeacherFullname
			object["status"] = res.Status.String()
			array = append(array, object)
		}
		result["reservations"] = array

		if data, err := json.Marshal(result); err == nil {
			fmt.Println(string(data))
			w.Write(data)
		}
	} else {

	}
}
