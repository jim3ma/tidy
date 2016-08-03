package system

import "gopkg.in/mgo.v2/bson"

type MailTemplate struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Type    string        `bson:"type" json:"type"`
	Content string        `bson:"content" json:"content"`
}
