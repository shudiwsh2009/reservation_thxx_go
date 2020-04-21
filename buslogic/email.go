package buslogic

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/mijia/sweb/log"
	"github.com/shudiwsh2009/reservation_thxx_go/config"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"github.com/shudiwsh2009/reservation_thxx_go/utils"
	"net/smtp"
	"strings"
)

func SendEmail(subject string, body string, receivers []string) error {
	for _, r := range receivers {
		if !utils.IsEmail(r) {
			return re.NewRErrorCode(fmt.Sprintf("wrong email %s", r), nil, re.ErrorFormatEmail)
		}
	}
	if config.Instance().IsStagingEnv() {
		log.Infof("SMOCK Send Email: \"%s\" to %s.\n", subject, strings.Join(receivers, ","))
		return nil
	}
	if len(receivers) == 0 {
		log.Info("empty email receivers")
		return nil
	}

	m := email.NewEmail()
	m.From = config.Instance().SMTPUser
	m.To = receivers
	m.Subject = subject
	m.Text = []byte(body)

	auth := smtp.PlainAuth("", config.Instance().SMTPUser, config.Instance().SMTPPassword, config.Instance().SMTPHost)
	if err := m.SendWithTLS(fmt.Sprintf("%s:%s", config.Instance().SMTPHost, config.Instance().SMTPPort),
		auth, &tls.Config{ServerName: config.Instance().SMTPHost}); err != nil {
		return re.NewRError(fmt.Sprintf("failed to send email %+v", m), err)
	}
	log.Infof("Send Email to %s, subject: \"%s\", body: \"%s\"", strings.Join(receivers, ","), subject, body)
	return nil
}

func EmailWarn(subject string, body string) error {
	return SendEmail(fmt.Sprintf("%s%s", "【thxxfzzx报警】", subject), body, config.Instance().EmailAddressDev)
}
