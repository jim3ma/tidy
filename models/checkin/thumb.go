package checkin

import "gopkg.in/mgo.v2/bson"

//"gopkg.in/mgo.v2"

// Thumb data
// the ID is the same as UserID
type Thumb struct {
	ID         bson.ObjectId   `bson:"_id" json:"id"`
	CheckinIDs []bson.ObjectId `bson:"cids" json:"cids"`
	//UserIDs []bson.ObjectId `bson:"uids" json:"-"`
}

func (t *Thumb) HasThumbed(cid bson.ObjectId) bool {
	for _, v := range t.CheckinIDs {
		if v == cid {
			return true
		}
	}
	return false
}
