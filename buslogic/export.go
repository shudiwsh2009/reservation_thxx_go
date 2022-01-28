package buslogic

import (
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"github.com/xuri/excelize/v2"
	"path/filepath"
	"sort"
	"time"
)

const ExportFileDefaultSheetName string = "export"

var StudentFeedbackChoice = map[uint8]string{
	'A': "非常同意",
	'B': "一般",
	'C': "不同意",
}

var ReservationRoomColumn = map[int][]string{
	1:  {"B", "C"},
	2:  {"D", "E"},
	3:  {"F", "G"},
	4:  {"H", "I"},
	5:  {"J", "K"},
	6:  {"L", "M"},
	7:  {"N", "O"},
	8:  {"P", "Q"},
	9:  {"R", "S"},
	10: {"T", "U"},
	11: {"V", "W"},
	12: {"X", "Y"},
}

func (w *Workflow) ExportReservationsToFile(reservations []*model.Reservation, path string) error {
	var err error
	file, err := excelize.OpenFile(filepath.Join(utils.ExportFolder, "export_template.xlsx"))
	if err != nil {
		return re.NewRError("fail to oepn export template", err)
	}
	for i, res := range reservations {
		rowNum := i + 2
		// 学生申请表
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("A%d", rowNum), res.StudentInfo.Fullname)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("B%d", rowNum), res.StudentInfo.Gender)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("C%d", rowNum), res.StudentInfo.Username)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("E%d", rowNum), res.StudentInfo.School)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("F%d", rowNum), res.StudentInfo.Hometown)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("G%d", rowNum), res.StudentInfo.Mobile)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("H%d", rowNum), res.StudentInfo.Email)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("J%d", rowNum), res.StudentInfo.Problem)
		// 预约信息
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("K%d", rowNum), res.TeacherFullname)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("L%d", rowNum), res.StartTime.Format("2006-01-02"))
		// 咨询师反馈表
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("O%d", rowNum), res.TeacherFeedback.Problem)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("P%d", rowNum), res.TeacherFeedback.Solution)
		// 学生反馈表
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("R%d", rowNum), res.StudentFeedback.Score)
		file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("S%d", rowNum), res.StudentFeedback.Feedback)
		if res.StudentFeedback.Choices != "" && len(res.StudentFeedback.Choices) == 12 {
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("T%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[0]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("U%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[1]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("V%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[2]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("W%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[3]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("X%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[4]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("Y%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[5]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("Z%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[6]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("AA%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[7]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("AB%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[8]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("AC%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[9]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("AD%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[10]])
			file.SetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("AE%d", rowNum), StudentFeedbackChoice[res.StudentFeedback.Choices[11]])
		}
	}
	err = file.SaveAs(path)
	if err != nil {
		return re.NewRError("fail to save to file", err)
	}
	return nil
}

// 咨询室
type ReservationRoom struct {
	Timetable map[string]*model.Reservation // key是开始时间，从08:00-18:00，每30分钟一档
}

// 咨询是否可以安排在本咨询室
func (room ReservationRoom) IsAvailableForReservation(reservation *model.Reservation) bool {
	start, end := reservation.GetStartAndEndTimeForArrangement()
	for ts := start; ts.Before(end); ts = ts.Add(30 * time.Minute) {
		if _, ok := room.Timetable[ts.Format("15:04")]; ok {
			return false
		}
	}
	return true
}

// 将咨询安排进咨询室
func (room ReservationRoom) AddReservation(reservation *model.Reservation) {
	start, end := reservation.GetStartAndEndTimeForArrangement()
	for ts := start; ts.Before(end); ts = ts.Add(30 * time.Minute) {
		room.Timetable[ts.Format("15:04")] = reservation
	}
}

// 安排表
type ReservationArrangement struct {
	Rooms []*ReservationRoom
}

func (w *Workflow) ExportReservationArrangementsToFile(reservations []*model.Reservation, path string) error {
	// 这是一个活动安排问题，贪心算法解决
	// 先按开始时间排序，这一步是贪心算法能获得最优解的关键
	sort.Sort(model.ByStartTimeOfReservation(reservations))
	// 建立映射表备用
	rId2ReservationMap := make(map[string]*model.Reservation)
	for _, r := range reservations {
		rId2ReservationMap[r.Id.Hex()] = r
	}
	// 初始化安排表
	arrangement := &ReservationArrangement{
		Rooms: make([]*ReservationRoom, 0),
	}
	for _, r := range reservations {
		// 寻找一个可用的咨询室
		hasArranged := false
		for _, room := range arrangement.Rooms {
			// 若存在就把咨询安排进去
			if room.IsAvailableForReservation(r) {
				room.AddReservation(r)
				hasArranged = true
				break
			}
		}
		// 找不到则申请一个新的咨询室
		if !hasArranged {
			newRoom := &ReservationRoom{
				Timetable: make(map[string]*model.Reservation),
			}
			newRoom.AddReservation(r)
			arrangement.Rooms = append(arrangement.Rooms, newRoom)
		}
	}

	// 开始写入文件
	var err error
	file := excelize.NewFile()
	file.SetDefaultFont("宋体")
	// 创建样式
	// 边框
	borderStyle, err := file.NewStyle(`{
		"border": [
		{
			"type": "left",
			"color": "000000",
			"style": 1
		},
		{
			"type": "top",
			"color": "000000",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "000000",
			"style": 1
		},
		{
			"type": "right",
			"color": "000000",
			"style": 1
		}]
	}`)
	if err != nil {
		return re.NewRError("failed to create border style", err)
	}
	// 文字居中带边框
	textCenterStyle, err := file.NewStyle(`{
		"alignment": {
			"horizontal": "center",
			"vertical": "center"
		},
		"border": [
			{
				"type": "left",
				"color": "000000",
				"style": 1
			},
			{
				"type": "top",
				"color": "000000",
				"style": 1
			},
			{
				"type": "bottom",
				"color": "000000",
				"style": 1
			},
			{
				"type": "right",
				"color": "000000",
				"style": 1
			}
		]
	}`)
	// 文字左上带边框
	textLeftTopStyle, err := file.NewStyle(`{
		"alignment": {
			"horizontal": "left",
			"vertical": "up"
		},
		"border": [
			{
				"type": "left",
				"color": "000000",
				"style": 1
			},
			{
				"type": "top",
				"color": "000000",
				"style": 1
			},
			{
				"type": "bottom",
				"color": "000000",
				"style": 1
			},
			{
				"type": "right",
				"color": "000000",
				"style": 1
			}
		]
	}`)

	// 咨询师排班表
	sheetIndex := file.NewSheet(ExportFileDefaultSheetName)
	file.SetActiveSheet(sheetIndex)
	// 第一表头
	file.SetCellValue(ExportFileDefaultSheetName, "A1", "时间")
	file.SetCellStyle(ExportFileDefaultSheetName, "A1", "A1", borderStyle)
	if len(arrangement.Rooms) > len(ReservationRoomColumn) {
		return re.NewRError("not enough reservation rooms", nil)
	}
	for i := 1; i <= len(arrangement.Rooms); i++ {
		cellIndex := fmt.Sprintf("%s1", ReservationRoomColumn[i][0])
		file.SetCellValue(ExportFileDefaultSheetName, cellIndex, fmt.Sprintf("咨询室%d", i))
		file.SetCellStyle(ExportFileDefaultSheetName, cellIndex, cellIndex, borderStyle)
	}
	// 合并第一表头
	for i := 1; i <= len(arrangement.Rooms); i++ {
		firstCellIndex := fmt.Sprintf("%s1", ReservationRoomColumn[i][0])
		secondCellIndex := fmt.Sprintf("%s1", ReservationRoomColumn[i][1])
		file.MergeCell(ExportFileDefaultSheetName, firstCellIndex, secondCellIndex)
		file.SetCellStyle(ExportFileDefaultSheetName, firstCellIndex, secondCellIndex, textCenterStyle)
	}
	// 第二表头
	for i := 1; i <= len(arrangement.Rooms); i++ {
		teacherCellIndex := fmt.Sprintf("%s2", ReservationRoomColumn[i][0])
		studentCellIndex := fmt.Sprintf("%s2", ReservationRoomColumn[i][1])
		file.SetCellValue(ExportFileDefaultSheetName, teacherCellIndex, "咨询师")
		file.SetCellStyle(ExportFileDefaultSheetName, teacherCellIndex, teacherCellIndex, borderStyle)
		file.SetCellValue(ExportFileDefaultSheetName, studentCellIndex, "学生")
		file.SetCellStyle(ExportFileDefaultSheetName, studentCellIndex, studentCellIndex, borderStyle)
	}
	// 合并第一列头
	file.MergeCell(ExportFileDefaultSheetName, "A1", "A2")
	file.SetCellStyle(ExportFileDefaultSheetName, "A1", "A1", textCenterStyle)
	// 填充咨询ID数据
	start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Local)
	end := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 18, 0, 0, 0, time.Local)
	for ts, rowIndex := start, 3; ts.Before(end); ts, rowIndex = ts.Add(30*time.Minute), rowIndex+1 {
		// 第一列的时间区间
		timeCellIndex := fmt.Sprintf("A%d", rowIndex)
		file.SetCellValue(ExportFileDefaultSheetName, timeCellIndex,
			fmt.Sprintf("%s-%s", ts.Format("15:04"), ts.Add(30*time.Minute).Format("15:04")))
		file.SetCellStyle(ExportFileDefaultSheetName, timeCellIndex, timeCellIndex, textCenterStyle)
		for i, room := range arrangement.Rooms {
			if r, ok := room.Timetable[ts.Format("15:04")]; ok {
				arrangementCellIndex := fmt.Sprintf("%s%d", ReservationRoomColumn[i+1][0], rowIndex)
				file.SetCellValue(ExportFileDefaultSheetName, arrangementCellIndex, r.Id.Hex())
				file.SetCellStyle(ExportFileDefaultSheetName, arrangementCellIndex, arrangementCellIndex, borderStyle)
			}
		}
	}
	// 合并咨询并填充详细信息
	timeRowNum := int(end.Sub(start) / (30 * time.Minute))
	for i := range arrangement.Rooms {
		for j := 3; j <= timeRowNum+2; {
			firstTeacherCellIndex := fmt.Sprintf("%s%d", ReservationRoomColumn[i+1][0], j)
			firstTeacherCellValue, _ := file.GetCellValue(ExportFileDefaultSheetName, firstTeacherCellIndex)
			if firstTeacherCellValue != "" {
				k := j
				for k <= timeRowNum+2 {
					kTeacherCellValue, _ := file.GetCellValue(ExportFileDefaultSheetName, fmt.Sprintf("%s%d", ReservationRoomColumn[i+1][0], k))
					if kTeacherCellValue != firstTeacherCellValue {
						break
					}
					k++
				}

				reservation, ok := rId2ReservationMap[firstTeacherCellValue]
				if !ok {
					log.Errorf("fail to find reservation %d", firstTeacherCellValue)
					continue
				}

				firstStudentCellIndex := fmt.Sprintf("%s%d", ReservationRoomColumn[i+1][1], j)
				file.SetCellValue(ExportFileDefaultSheetName, firstTeacherCellIndex,
					fmt.Sprintf("%s %s %s", reservation.TeacherFullname, reservation.TeacherMobile, reservation.TeacherAddress))
				file.SetCellStyle(ExportFileDefaultSheetName, firstTeacherCellIndex, firstStudentCellIndex, textLeftTopStyle)
				file.SetCellValue(ExportFileDefaultSheetName, firstStudentCellIndex,
					fmt.Sprintf("%s %s", reservation.StudentInfo.Fullname, reservation.StudentInfo.Mobile))
				file.SetCellStyle(ExportFileDefaultSheetName, firstStudentCellIndex, firstStudentCellIndex, textLeftTopStyle)
				for mergeRow := j + 1; mergeRow < k; mergeRow++ {
					mergeTeacherCellIndex := fmt.Sprintf("%s%d", ReservationRoomColumn[i+1][0], mergeRow)
					mergeStudentCellIndex := fmt.Sprintf("%s%d", ReservationRoomColumn[i+1][1], mergeRow)
					file.MergeCell(ExportFileDefaultSheetName, firstTeacherCellIndex, mergeTeacherCellIndex)
					file.SetCellStyle(ExportFileDefaultSheetName, firstTeacherCellIndex, mergeTeacherCellIndex, textLeftTopStyle)
					file.MergeCell(ExportFileDefaultSheetName, firstStudentCellIndex, mergeStudentCellIndex)
					file.SetCellStyle(ExportFileDefaultSheetName, firstStudentCellIndex, mergeStudentCellIndex, textLeftTopStyle)
				}
				j = k
			} else {
				j++
			}
		}
	}

	file.SetActiveSheet(sheetIndex)
	err = file.SaveAs(path)
	if err != nil {
		return re.NewRError("fail to save file to path", err)
	}
	return nil
}
