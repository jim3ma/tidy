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
	Setting     Setting       `bson:"setting" json:"setting"`
}

func (u *User) CanCheckIn() bool {
	if u.LastCheckIn == nil {
		return true
	}
	//log.Printf("Last checkin status: %s, current time: %s", u.LastCheckIn, time.Now())
	ciData := new(ci.CheckIn)
	if checkin, ok := u.LastCheckIn.(bson.M); ok {
		if cb, err := bson.Marshal(checkin); err == nil {
			if err = bson.Unmarshal(cb, ciData); err == nil {
				now := time.Now()
				if ciData.CreateYear == now.Year() &&
					ciData.CreateMonth == int(now.Month()) &&
					ciData.CreateDay == now.Day() {
					return false
				}
			}
		}
	}
	return true
}
