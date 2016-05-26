package utilities

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

type MailConfig struct {
	AuthMailAddr  string
	AuthPassword  string
	MailFrom      string
	SmtpHost      string
	TlsSkipVerify bool
}

var config MailConfig

// StartTLS Email
func InitConfig(c MailConfig) {
	config = c
}

func SendSysMail(mailto string, subject string, body string) error {

	from := mail.Address{"", config.MailFrom}
	to := mail.Address{"", mailto}
	//subject := "This is the email subject"
	body := "This is a body.\n With two lines."

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	host, _, _ := net.SplitHostPort(config.SmtpHost)

	auth := smtp.PlainAuth("", config.AuthMailAddr, config.AuthPassword, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: config.TlsSkipVerify,
		ServerName:         host,
	}

	c, err := smtp.Dial(config.SmtpHost)
	if err != nil {
		log.Panic(err)
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}
