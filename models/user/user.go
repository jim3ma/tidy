package user

import (
	"time"

	ci "github.com/jim3mar/tidy/models/checkin"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"uid"`
	UserName    string        `bson:"user_name" json:"user_name"`
	Password    string        `bson:"password" json:"-"`
	EMail       string        `bson:"email" json:"email"`
	CreateAt    time.Time     `bson:"create_at" json:"-"`
	Timestamp   int64         `bson:"timestamp" json:"-"`
	Portrait    string        `bson:"portrait" json:"portrait"`
	Continuous  int           `bson:"continuous" json:"continuous"`
	LastCheckIn interface{}   `bson:"last_checkin" json:"last_checkin"`
}

func (u *User) CanCheckIn() bool {
	if u.LastCheckIn == nil {
		return false
	}
	checkin := u.LastCheckIn.(ci.CheckIn)
	now := time.Now()
	if checkin.CreateYear == now.Year() &&
		checkin.CreateMonth == int(now.Month()) &&
		checkin.CreateDay == now.Day() {
		return true
	}
	return false
}
