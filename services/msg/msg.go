package msg

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type MsgResource struct {
	Mongo   *mgo.Session
	CollMsg *mgo.Collection
}

func (mr *MsgResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	mr.Mongo = session
	mr.CollMsg = mr.Mongo.DB(db).C("message")
}

func (mr *MsgResource) QueryMsgByUserId(c *gin.Context) {

}
