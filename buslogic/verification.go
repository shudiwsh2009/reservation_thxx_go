package buslogic

type StudentVerification struct {
	username string
	fullname string
	verified bool // 是否有学籍
}

var username2StudentVerification = make(map[string]*StudentVerification)

func (w *Workflow) loadStudentVerification(filename string) {

}

func (w *Workflow) verifyStudent(username string, fullname string) bool {
	verification, ok := username2StudentVerification[username]
	return ok && verification.fullname == fullname && verification.verified
}
