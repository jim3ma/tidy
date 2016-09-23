package checkin

import (
	"github.com/gin-gonic/gin"
	mod "github.com/jim3ma/tidy/models/checkin"
	"gopkg.in/mgo.v2/bson"
)

func (cr *CheckInResource) ListFavorByUserID(c *gin.Context) {
	uid := bson.ObjectIdHex(c.Query("uid"))
	var favor mod.Favorite
	cr.CollFavor.Find(
		bson.M{
			"_id": uid,
		},
	).One(&favor)
}
