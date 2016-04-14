package user

import (
	//"github.com/jim3mar/tidy/models/checkin
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id_       bson.ObjectId `bson:"_id" json:"uid"`
	UserName  string        `bson:"username" json:"username"`
	Password  string        `bson:"password" json:"password"`
	EMail     string        `bson:"email" json:"email"`
	CreateAt  time.Time     `bson:"create_at" json:"create_at"`
	Timestamp int64         `bson:"timestamp" json:"timestamp"`
	Portrait  string        `bson:"portrait" json:"portrait"`
        Continuous int          `bson:"continuous" json:"continuous"`
        LastCheckIn interface{} `bson:"last_checkin" json:last_checkin`
}

