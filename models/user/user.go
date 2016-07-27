package user

import (
	"time"

	ci "github.com/jim3mar/tidy/models/checkin"
	"github.com/spf13/viper"
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
	if viper.GetBool("debug") {
		return true
	}
	if u.LastCheckIn == nil {
		return true
	}
	ciData := new(ci.CheckIn)
	if checkin, ok := u.LastCheckIn.(bson.M); ok {
		if bson2Struct(&checkin, ciData) {
			now := time.Now()
			if ciData.CreateYear == now.Year() &&
				ciData.CreateMonth == int(now.Month()) &&
				ciData.CreateDay == now.Day() {
				return false
			}
		}
	}
	return true
}

// TBD
func (u *User) CheckInStatus() (can bool, cont int) {
	if u.LastCheckIn == nil {
		can = true
		cont = 1
		return
	}
	//log.Infof("Last checkin status: %s, current time: %s", u.LastCheckIn, time.Now())
	ciData := new(ci.CheckIn)
	if checkin, ok := u.LastCheckIn.(bson.M); ok {
		if bson2Struct(&checkin, ciData) {
			now := time.Now()
			if ciData.CreateYear == now.Year() &&
				ciData.CreateAt.YearDay()+1 == now.YearDay() {
				cont = u.Continuous + 1
				can = true
				return
			} else if ciData.CreateYear == now.Year() &&
				ciData.CreateAt.YearDay() == now.YearDay() {
				cont = 1
				can = false
				return
			}
		}
	}
	can = true
	cont = 1
	return
}

func (u *User) CalcContinuous() int {
	if u.LastCheckIn == nil {
		return 1
	}
	ciData := new(ci.CheckIn)
	if checkin, ok := u.LastCheckIn.(bson.M); ok {
		if bson2Struct(&checkin, ciData) {
			now := time.Now()
			if ciData.CreateYear == now.Year() &&
				ciData.CreateAt.YearDay()+1 == now.YearDay() {
				return u.Continuous + 1
			}
		}
	}
	return 1
}

func bson2Struct(bs *bson.M, st interface{}) bool {
	if cb, err := bson.Marshal(bs); err == nil {
		if err = bson.Unmarshal(cb, st); err == nil {
			return true
		}
	}
	return false
}
