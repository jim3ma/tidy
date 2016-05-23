package system

import (
    "github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type SystemResource struct {
	Mongo    *mgo.Session
	CollSys  *mgo.Collection
    CollFeedback *mgo.Collection
}

func (sr *SystemResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	sr.Mongo = session
	sr.CollSys = sr.Mongo.DB(db).C("system")
	sr.CollFeedback = sr.Mongo.DB(db).C("feedback")
}
