package buslogic

import (
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"github.com/tealeg/xlsx"
)

type StudentVerification struct {
	username string
	fullname string
	verified bool // 是否有学籍
}

var username2StudentVerification = make(map[string]*StudentVerification)

func (w *Workflow) loadStudentVerification(filename string) error {
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		return err
	}
	sheet := file.Sheets[0]
	for _, row := range sheet.Rows {
		if len(row.Cells) != 3 {
			break
		}
		username := row.Cells[0].Value
		if !utils.IsStudentUsername(username) {
			continue
		}
		fullname := row.Cells[1].Value
		verified := row.Cells[2].Value
		username2StudentVerification[username] = &StudentVerification{
			username: username,
			fullname: fullname,
			verified: verified == "是",
		}
	}
	return nil
}

func (w *Workflow) verifyStudent(username string, fullname string) bool {
	verification, ok := username2StudentVerification[username]
	return ok && verification.fullname == fullname && verification.verified
}
