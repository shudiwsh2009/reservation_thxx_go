package buslogic

import (
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/model"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"sort"
	"time"
)

// 定义单元格样式
var (
	textCenterAlignment = xlsx.Alignment{
		Horizontal: "center",
		Vertical:   "center",
	}
	textTopAlignment = xlsx.Alignment{
		Horizontal: "left",
		Vertical:   "up",
	}
	border = xlsx.Border{
		Left:        "thin",
		LeftColor:   "000000",
		Right:       "thin",
		RightColor:  "000000",
		Top:         "thin",
		TopColor:    "000000",
		Bottom:      "thin",
		BottomColor: "000000",
	}

	borderStyle *xlsx.Style = &xlsx.Style{
		ApplyBorder: true,
		Border:      border,
	}
	textCenterStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyBorder:    true,
		Alignment:      textCenterAlignment,
		Border:         border,
	}
	textTopStyle *xlsx.Style = &xlsx.Style{
		ApplyAlignment: true,
		ApplyBorder:    true,
		Alignment:      textTopAlignment,
		Border:         border,
	}
)

func (w *Workflow) ExportReservationsToFile(reservations []*model.Reservation, path string) error {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	file, err = xlsx.OpenFile(filepath.Join(utils.ExportFolder, "export_template.xlsx"))
	if err != nil {
		return re.NewRError("fail to oepn export template", err)
	}
	sheet = file.Sheet["export"]
	if sheet == nil {
		return re.NewRError("fail to open sheet", err)
	}
	for _, res := range reservations {
		row = sheet.AddRow()
		// 学生申请表
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Fullname)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Gender)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Username)
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.School)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Hometown)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Mobile)
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Email)
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.StudentInfo.Problem)
		// 预约信息
		cell = row.AddCell()
		cell.SetString(res.TeacherFullname)
		cell = row.AddCell()
		cell.SetString(res.StartTime.Format("2006-01-02"))
		// 咨询师反馈表
		cell = row.AddCell()
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.TeacherFeedback.Problem)
		cell = row.AddCell()
		cell.SetString(res.TeacherFeedback.Solution)
		// 学生反馈表
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(res.StudentFeedback.Score)
		cell = row.AddCell()
		cell.SetString(res.StudentFeedback.Feedback)
		if res.StudentFeedback.Choices != "" {
			for i := 0; i < len(res.StudentFeedback.Choices); i++ {
				cell = row.AddCell()
				switch res.StudentFeedback.Choices[i] {
				case 'A':
					cell.SetString("非常同意")
				case 'B':
					cell.SetString("一般")
				case 'C':
					cell.SetString("不同意")
				default:
				}
			}
		}
	}
	err = file.Save(path)
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
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	xlsx.SetDefaultFont(15, "宋体")
	file = xlsx.NewFile()
	// 咨询师排班表
	sheet, err = file.AddSheet("咨询师排班表")
	if err != nil {
		return re.NewRError("fail to add sheet for arrangement", err)
	}
	// 第一表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	cell.SetValue("时间")
	for i := 1; i <= len(arrangement.Rooms); i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue(fmt.Sprintf("咨询室%d", i))
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
	}
	// 合并第一表头
	for i := 0; i < len(arrangement.Rooms); i++ {
		cell = row.Cells[2*i+1]
		cell.Merge(1, 0)
		cell.SetStyle(textCenterStyle)
	}
	// 第二表头
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(borderStyle)
	for i := 0; i < len(arrangement.Rooms); i++ {
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue("咨询师")
		cell = row.AddCell()
		cell.SetStyle(borderStyle)
		cell.SetValue("学生")
	}
	// 合并第一列头
	row = sheet.Rows[0]
	cell = row.Cells[0]
	cell.Merge(0, 1)
	cell.SetStyle(textCenterStyle)
	// 填充咨询ID数据
	start := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Local)
	end := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 18, 0, 0, 0, time.Local)
	for ts := start; ts.Before(end); ts = ts.Add(30 * time.Minute) {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(textCenterStyle)
		cell.SetValue(fmt.Sprintf("%s-%s", ts.Format("15:04"), ts.Add(30 * time.Minute).Format("15:04")))
		for _, room := range arrangement.Rooms {
			if r, ok := room.Timetable[ts.Format("15:04")]; ok {
				cell = row.AddCell()
				cell.SetStyle(borderStyle)
				cell.SetValue(r.Id.Hex())
				cell = row.AddCell()
				cell.SetStyle(borderStyle)
			} else {
				cell = row.AddCell()
				cell.SetStyle(borderStyle)
				cell = row.AddCell()
				cell.SetStyle(borderStyle)
			}
		}
	}
	// 合并咨询并填充详细信息
	timeRowNum := int(end.Sub(start) / (30 * time.Minute))
	for i := range arrangement.Rooms {
		for j := 2; j < timeRowNum+2; {
			firstTeacherCell := sheet.Cell(j, 2*i+1)
			if firstTeacherCell.Value != "" {
				k := j
				for k < timeRowNum+2 {
					kTeacherCell := sheet.Cell(k, 2*i+1)
					if kTeacherCell.Value != firstTeacherCell.Value {
						break
					}
					k++
				}
				reservation, ok := rId2ReservationMap[firstTeacherCell.Value]
				if !ok {
					log.Errorf("fail to find reservation %d", firstTeacherCell.Value)
					continue
				}
				firstStudentCell := sheet.Cell(j, 2*i+2)
				firstTeacherCell.Merge(0, k-j-1)
				firstTeacherCell.SetValue(textTopStyle)
				firstTeacherCell.SetValue(fmt.Sprintf("%s\n%s\n%s",
					reservation.TeacherFullname, reservation.TeacherMobile, reservation.TeacherAddress))
				firstStudentCell.Merge(0, k-j-1)
				firstStudentCell.SetValue(textTopStyle)
				firstStudentCell.SetValue(fmt.Sprintf("%s\n%s",
					reservation.StudentInfo.Fullname, reservation.StudentInfo.Mobile))
				j = k
			} else {
				j++
			}
		}
	}
	// 调整列宽
	sheet.SetColWidth(0, 0, 12)
	sheet.SetColWidth(1, 2*len(arrangement.Rooms), 20)

	err = file.Save(path)
	if err != nil {
		return re.NewRError("fail to save file to path", err)
	}
	return nil
}
