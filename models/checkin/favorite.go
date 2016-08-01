package checkin

import "gopkg.in/mgo.v2/bson"

//"gopkg.in/mgo.v2"

// Favorite data
// the ID is the same as UserID
type Favorite struct {
	ID         bson.ObjectId   `bson:"_id" json:"id"`
	CheckinIDs []bson.ObjectId `bson:"cids" json:"cids"`
	//UserID     bson.ObjectId   `bson:"uid" json:"uid"`
}

func (f *Favorite) HasFavored(cid bson.ObjectId) bool {
	for _, v := range f.CheckinIDs {
		if v == cid {
			return true
		}
	}
	return false
}
