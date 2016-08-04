package system

import (
	"bytes"
	"net/mail"
	"text/template"

	log "github.com/Sirupsen/logrus"
	modsys "github.com/jim3mar/tidy/models/system"
	mod "github.com/jim3mar/tidy/models/user"
	util "github.com/jim3mar/tidy/utilities"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SysResource contains mail operation
type SysResource struct {
	Mongo            *mgo.Session
	CollSys          *mgo.Collection
	CollFeedback     *mgo.Collection
	CollMailTemplate *mgo.Collection
}

// Init config env
func (sr *SysResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	sr.Mongo = session
	sr.CollSys = sr.Mongo.DB(db).C("system")
	sr.CollFeedback = sr.Mongo.DB(db).C("feedback")
	sr.CollMailTemplate = sr.Mongo.DB(db).C("sys_mail_template")
}

// Mail type
const (
	MTWelcome = iota
	MTResetPWD
	MTSysInfo
)

func (sr *SysResource) sendMail(mailType int, user *mod.User, data interface{}) bool {
	// TBD
	// cache
	var mailtmpl modsys.MailTemplate
	err := sr.CollMailTemplate.Find(bson.M{
		"type": mailType,
	}).Sort("-version").One(&mailtmpl)
	if err != nil {
		log.Errorf("Query mail template failed, error: %s", err)
		return false
	}
	content, err := mailtmpl.GenerateContent(data)
	if err != nil {
		// TBD
		return false
	}
	subject, err := mailtmpl.GenerateSubject(data)
	if err != nil {
		// TBD
		return false
	}
	mailto := mail.Address{
		Name:    user.UserName,
		Address: user.EMail,
	}
	err = util.SendSysMail(mailto, subject, content)
	if err != nil {
		return false
	}
	return true
}

// SendResetPWDMail send email for resetting password
func (sr *SysResource) SendResetPWDMail(user *mod.User, authToken string) bool {
	subject := "重置密码"
	body := `<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=2.0">
    </head>
    <body>
	    Hi, {{ .User.UserName }}
		<br />
	    <a href="http://tf.ctidy.com/user/password.html?auth_token={{- .AuthToken -}}">点击这里重置密码</a>
        <br />
    </body>
</html>`
	type Body struct {
		User      mod.User
		AuthToken string
	}
	tmpl := template.Must(template.New("resetpwd").Parse(body))
	var buf bytes.Buffer
	tmpl.Execute(&buf, Body{
		User:      *user,
		AuthToken: authToken,
	})

	log.Debugf("Send resetting password email to %s, auth_token: %s", user.EMail, authToken)
	log.Debugf("Content %s", buf.String())
	mailto := mail.Address{
		Name:    user.UserName,
		Address: user.EMail,
	}
	util.SendSysMail(mailto, subject, buf.String())
	// util.SendSysMail(mailto, subject, body)
	return false
}

// SendWelcomeMail send email to user after registered
func (sr *SysResource) SendWelcomeMail(user *mod.User) bool {
	subject := "{{ .UserName }}, 欢迎来到Tidy"
	body := `<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=2.0">
    </head>
    <body>
	    <h2>{{ .UserName }}欢迎来到Tidy</h2>
        <br/>
    </body>
</html>`

	tmpl := template.Must(template.New("welcome").Parse(body))
	var buf bytes.Buffer
	tmpl.Execute(&buf, *user)

	log.Debugf("Send welcome email to %+v", *user)
	log.Debugf("Content %s", buf.String())
	mailto := mail.Address{
		Name:    user.UserName,
		Address: user.EMail,
	}
	util.SendSysMail(mailto, subject, buf.String())
	// util.SendSysMail(mailto, subject, body)
	return false
}
