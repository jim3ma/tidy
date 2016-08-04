package system

import (
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2"
	//"time"
)

// Feedback is used for system bug report
type Feedback struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	UserName  string        `bson:"username" json:"user_name"`
	Content   string        `bson:"content" json:"content"`
	Timestamp int64         `bson:"timestamp" json:"timestamp"`
}
