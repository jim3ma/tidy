package checkin

import (
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2"
	"time"
)

type CheckIn struct {
	Id_         bson.ObjectId `bson:"_id" json:"id"`
	UserId      bson.ObjectId `bson:"user_id" json:"user_id"`
	Content     string        `bson:"content" json:"content"`
	CreateAt    time.Time     `bson:"create_at" json:"create_at"`
	CreateDay   int           `bson:"create_day" json:"create_day"`
	CreateMonth int           `bson:"create_month" json:"create_month"`
	CreateYear  int           `bson:"create_year" json:"create_year"`
	CreateHour  int           `bson:"create_hour" json:"create_hour"`
	CreateMin   int           `bson:"create_min" json:"create_min"`
	CreateSec   int           `bson:"create_sec" json:"create_sec"`
	Timestamp   int64         `bson:"timestamp" json:"timestamp"`
	Images      []string      `bson:"images" json:"images"`
}
