package utils

import (
	"errors"
	"fmt"
	"github.com/shudiwsh2009/reservation_thxx_go/models"
	"github.com/tealeg/xlsx"
	"strings"
)

const (
	DefaultReservationExportExcelFilename = "export_template.xlsx"
	ExportFolder                          = "assets/export/"
	ExcelSuffix                           = ".xlsx"
)

func ExportReservationsToExcel(reservations []*models.Reservation, filename string) error {
	xl, err := xlsx.OpenFile(ExportFolder + DefaultReservationExportExcelFilename)
	if err != nil {
		return errors.New("导出失败:打开模板文件失败")
	}
	sheet := xl.Sheet["export"]
	if sheet == nil {
		return errors.New("导出失败:打开工作表失败")
	}
	var row *xlsx.Row
	var cell *xlsx.Cell
	for _, reservation := range reservations {
		fmt.Println(reservation)
		row = sheet.AddRow()
		// 学生申请表
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.Name)
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.Gender)
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.StudentId)
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.School)
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.Hometown)
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.Mobile)
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.Email)
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(reservation.StudentInfo.Problem)
		// 预约信息
		cell = row.AddCell()
		cell.SetString(reservation.TeacherFullname)
		cell = row.AddCell()
		cell.SetString(reservation.StartTime.Format("2006-01-02"))
		// 咨询师反馈表
		cell = row.AddCell()
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(reservation.TeacherFeedback.Problem)
		cell = row.AddCell()
		cell.SetString(reservation.TeacherFeedback.Solution)
		// 学生反馈表
		cell = row.AddCell()
		cell = row.AddCell()
		cell.SetString(reservation.StudentFeedback.Score)
		cell = row.AddCell()
		cell.SetString(reservation.StudentFeedback.Feedback)
		if !strings.EqualFold(reservation.StudentFeedback.Choices, "") {
			for i := 0; i < len(reservation.StudentFeedback.Choices); i++ {
				cell = row.AddCell()
				switch reservation.StudentFeedback.Choices[i] {
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
	err = xl.Save(ExportFolder + filename)
	if err != nil {
		return errors.New("导出失败:保存文件失败")
	}
	return nil
}
