package user

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id_      bson.ObjectId `bson:"_id"`
	UserName string        `bson:"username"`
	Password string        `bson:"password"`
	EMail    string        `bson:"email"`
	CreateAt time.Time     `bson:"create_at"`
	Portrait string        `bson:"portrait"`
}
