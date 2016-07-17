package msg

import (
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2"
	//"time"
)

type Msg struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	From      bson.ObjectId `bson:"from" json:"from"`
	To        bson.ObjectId `bson:"to" json:"to"`
	Content   string        `bson:"content" json:"content"`
	Timestamp int64         `bson:"timestamp" json:"timestamp"`
}
