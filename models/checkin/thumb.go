package checkin

import "gopkg.in/mgo.v2/bson"

//"gopkg.in/mgo.v2"

// Thumb data
// the ID is the same as UserID
type Thumb struct {
	ID         bson.ObjectId   `bson:"_id" json:"id"`
	CheckinIDs []bson.ObjectId `bson:"cids" json:"cids"`
	//Count   int             `bson:"count" json:"count"`
	//UserIDs []bson.ObjectId `bson:"uids" json:"-"`
}
