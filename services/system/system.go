package system

import (
	log "github.com/Sirupsen/logrus"
	mod "github.com/jim3mar/tidy/models/user"
	util "github.com/jim3mar/tidy/utilities"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

type SystemResource struct {
	Mongo        *mgo.Session
	CollSys      *mgo.Collection
	CollFeedback *mgo.Collection
}

func (sr *SystemResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	sr.Mongo = session
	sr.CollSys = sr.Mongo.DB(db).C("system")
	sr.CollFeedback = sr.Mongo.DB(db).C("feedback")
}

func (sr *SystemResource) SendResetPWDMail(user *mod.User, authToken string) bool {
	subject := "重置密码"
	body := `<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=2.0">
    </head>
    <body>
	    <a href="http://tf.ctidy.com/user/password.html?auth_token=` + authToken + `">点击这里重置密码</a>
        <br/>
    </body>
</html>`

	log.Debugf("Send resetting password email to %s, auth_token: %s", user.EMail, authToken)
	util.SendSysMail(user.EMail, subject, body)
	return false
}
