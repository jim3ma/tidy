package checkin

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	mod "github.com/jim3ma/tidy/models/checkin"
	"gopkg.in/mgo.v2/bson"
)

func (cr *CheckInResource) ListThumbByUserID(c *gin.Context) {
	uid := bson.ObjectIdHex(c.Query("uid"))
	var thumb mod.Thumb
	cr.CollThumb.Find(
		bson.M{
			"_id": uid,
		},
	).One(&thumb)
	l := len(thumb.CheckinIDs)
	cis := make([]mod.CheckIn, l, l)
	for i, cid := range thumb.CheckinIDs {
		ci, err := cr.queryCheckInByCID(cid)
		if err == nil {
			cis[i] = ci
		} else {
			log.Errorf("query checkin error: %s", err)
		}
	}
	// TBD
	// update the queried ci
}
