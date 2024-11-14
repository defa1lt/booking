package email

import "gopkg.in/gomail.v2"

type EmailSender struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

func NewEmailSender(from, host string, port int, username, password string) *EmailSender {
	return &EmailSender{
		From:     from,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (e *EmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)

	return d.DialAndSend(m)
}
