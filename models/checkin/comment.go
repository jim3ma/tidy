package checkin

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//"gopkg.in/mgo.v2"

// Comment data
// the ID is the same as CheckinID
type Comment struct {
	ID         bson.ObjectId   `bson:"_id" json:"id"`
	CommentIDs []bson.ObjectId `bson:"comment_ids" json:"comment_ids"`
	//CheckinID bson.ObjectId   `bson:"checkin_id" json:"checkin_id"`
}

// SingleComment for Comment data
type SingleComment struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	UserID      bson.ObjectId `bson:"uid" json:"uid"`
	UserName    string        `bson:"user_name" json:"user_name"`
	Content     string        `bson:"content" json:"content"`
	CreateAt    time.Time     `bson:"create_at" json:"create_at"`
	CreateDay   int           `bson:"create_day" json:"create_day"`
	CreateMonth int           `bson:"create_month" json:"create_month"`
	CreateYear  int           `bson:"create_year" json:"create_year"`
	CreateHour  int           `bson:"create_hour" json:"create_hour"`
	CreateMin   int           `bson:"create_min" json:"create_min"`
	CreateSec   int           `bson:"create_sec" json:"create_sec"`
	Timestamp   int64         `bson:"timestamp" json:"timestamp"`
}
