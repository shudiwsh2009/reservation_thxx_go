package buslogic

import (
	"fmt"
	"github.com/mijia/sweb/log"
	"github.com/scorredoira/email"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"net/mail"
	"net/smtp"
	"strings"
)

func SendEmail(m *email.Message) error {
	if config.Instance().IsStagingEnv() {
		log.Infof("SMOCK Send Email: \"%s\" to %s.\n", m.Subject, strings.Join(config.Instance().EmailAddressDev, ","))
		return nil
	}

	auth := smtp.PlainAuth("", config.Instance().SMTPUser, config.Instance().SMTPPassword, config.Instance().SMTPHost)
	if err := email.Send(config.Instance().SMTPHost+":587", auth, m); err != nil {
		return re.NewRError(fmt.Sprintf("failed to send email %+v", m), err)
	}
	return nil
}

func EmailWarn(subject string, body string) error {
	if len(config.Instance().EmailAddressDev) == 0 {
		return nil
	}
	m := email.NewMessage(subject, body)
	m.From = mail.Address{Name: "", Address: config.Instance().SMTPUser}
	m.To = config.Instance().EmailAddressDev
	return SendEmail(m)
}
