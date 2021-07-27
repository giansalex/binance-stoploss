package notify

import (
	"strings"

	"gopkg.in/gomail.v2"
)

type MailNotify struct {
	sender  gomail.SendCloser
	subject string
	from    string
	to      string
}

func NewMailNotify(host string, port int, username, password, subject, from, to string) *MailNotify {
	d := gomail.NewDialer(host, port, username, password)
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	return &MailNotify{s, subject, from, to}
}

// Send notify to email
func (mailNotify *MailNotify) Send(message string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailNotify.from)
	m.SetHeader("To", mailNotify.to)
	m.SetHeader("Subject", mailNotify.subject)
	m.SetBody("text/html", mailNotify.fixLineEnding(message))

	return gomail.Send(mailNotify.sender, m)
}

func (mailNotify *MailNotify) fixLineEnding(msg string) string {
	return strings.ReplaceAll(msg, "\n", "<br/>")
}
