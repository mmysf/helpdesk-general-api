package helpers

import (
	"app/domain/model"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	From(fromEmail, fromName string)
	To(receiver []string)
	Subject(value string)
	Body(value string)
	Attachment(r io.Reader, filename string, c string)
	AttachmentFile(filename string)
	Send() error
}

type smtpMailer struct {
	email  *gomail.Message
	dialer *gomail.Dialer
}

func (mailer *smtpMailer) From(fromEmail, fromName string) {
	mailer.email.SetHeader("From", fmt.Sprintf("%s <%s>", fromName, fromEmail))
}

func (mailer *smtpMailer) To(val []string) {
	mailer.email.SetHeader("To", val...)
}

func (mailer *smtpMailer) Subject(val string) {
	mailer.email.SetHeader("Subject", val)
}

func (mailer *smtpMailer) Body(val string) {
	mailer.email.SetBody("text/html", val)
}

func (mailer *smtpMailer) Attachment(r io.Reader, filename string, c string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	mailer.email.Attach(
		filename,
		gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(buf.Bytes())
			return err
		}),
		gomail.SetHeader(map[string][]string{"Content-Type": {c}}),
	)
}

func (mailer *smtpMailer) AttachmentFile(filename string) {
	mailer.email.Attach(filename)
}

func (mailer *smtpMailer) Send() error {
	return mailer.dialer.DialAndSend(mailer.email)
}

func NewSMTPMailer(company *model.Company) Mailer {
	// init mail
	mail := gomail.NewMessage()

	// d := gomail.NewDialer("smtp.example.com", 587, "user", "123456")
	mailport, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	d := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		mailport,
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"),
	)

	// If the SMTP server requires a TLS connection, you can set it like this
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// For port 25, set SSL = false
	// For port 465, set SSL = true (SMTPS)
	// For port 587, set SSL = true (with STARTTLS)

	mailer := smtpMailer{
		email:  mail,
		dialer: d,
	}

	// default sender
	mailer.From(company.Settings.SMTP.FromAddress, company.Settings.SMTP.FromName)

	return &mailer
}
