package user

import (
	//"github.com/jim3mar/tidy/models/checkin
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id_         bson.ObjectId `bson:"_id" json:"uid"`
	UserName    string        `bson:"username" json:"username"`
	Password    string        `bson:"password" json:"-"`
	EMail       string        `bson:"email" json:"email"`
	CreateAt    time.Time     `bson:"create_at" json:"-"`
	Timestamp   int64         `bson:"timestamp" json:"-"`
	Portrait    string        `bson:"portrait" json:"portrait"`
	Continuous  int           `bson:"continuous" json:"continuous"`
	LastCheckIn interface{}   `bson:"last_checkin" json:last_checkin`
}
