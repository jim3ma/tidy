package utilities

import (
	"crypto/tls"
	"fmt"
	//"log"
	"net"
	"net/mail"
	"net/smtp"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// MailConfig contains all configuration for mail
type MailConfig struct {
	AuthMailAddr  string
	AuthPassword  string
	SendFrom      string
	SMTPHost      string
	TLSSkipVerify bool
}

var config MailConfig

// TLS Email

// InitMailConfig setup mail configuration
func InitMailConfig() {
	config.AuthMailAddr = viper.GetString("mail.authaddr")
	config.AuthPassword = viper.GetString("mail.authpasswd")
	config.SMTPHost = fmt.Sprintf("%s:%s",
		viper.GetString("mail.host"), viper.GetString("mail.port"))
	config.SendFrom = viper.GetString("mail.sendfrom")
	config.TLSSkipVerify = viper.GetBool("mail.tlsskipverify")
	log.Infof("current mail config: %+v", config)
}

//SendSysMail use global variable config for default smtp settings.
//just put mailto, subject, and body
func SendSysMail(mailto string, subject string, body string) error {

	from := mail.Address{
		Name:    "",
		Address: config.SendFrom,
	}
	to := mail.Address{
		Name:    "",
		Address: mailto,
	}
	//subject := "This is the email subject"
	//body := "This is a body.\n With two lines."

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
	host, _, _ := net.SplitHostPort(config.SMTPHost)

	auth := smtp.PlainAuth("", config.AuthMailAddr, config.AuthPassword, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: config.TLSSkipVerify,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", config.SMTPHost, tlsconfig)
	if err != nil {
		log.Errorf("Cann't connect to smtp server, err: %s", err)
		return err
	}
	//////////////////////////////////////
	//c, err := smtp.Dial(config.SMTPHost)
	//if err != nil {
	//	log.Panic(err)
	//}

	//c.StartTLS(tlsconfig)
	//////////////////////////////////////

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	log.Infof("Send mail to: %s", mailto)
	log.Infof("Subject: %s", subject)
	log.Infof("Content: %s", body)
	return nil
}

// SendSysMailWithStartTLS use global variable config for default smtp settings.
// just put mailto, subject, and body
func SendSysMailWithStartTLS(mailto string, subject string, body string) error {

	from := mail.Address{
		Name:    "",
		Address: config.SendFrom,
	}
	to := mail.Address{
		Name:    "",
		Address: mailto,
	}
	//subject := "This is the email subject"
	//body := "This is a body.\n With two lines."

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
	host, _, _ := net.SplitHostPort(config.SMTPHost)

	auth := smtp.PlainAuth("", config.AuthMailAddr, config.AuthPassword, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: config.TLSSkipVerify,
		ServerName:         host,
	}

	c, err := smtp.Dial(config.SMTPHost)
	if err != nil {
		log.Panic(err)
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}
