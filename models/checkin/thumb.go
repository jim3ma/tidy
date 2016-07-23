package checkin

import "gopkg.in/mgo.v2/bson"

//"gopkg.in/mgo.v2"

// Thumb data
// the ID is the same as CheckinID
type Thumb struct {
	ID      bson.ObjectId   `bson:"_id" json:"id"`
	Count   int             `bson:"count" json:"count"`
	UserIDs []bson.ObjectId `bson:"uids" json:"-"`
	//CheckinID bson.ObjectId   `bson:"cid" json:"cid"`
}
