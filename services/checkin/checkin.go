package checkin

import (
	"log"

	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/checkin"
	//"github.com/jim3mar/tidy/models/user"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"encoding/json"
	//"log"
	//"strconv"
	"time"
)

type CheckInResource struct {
	Mongo    *mgo.Session
	CollCI   *mgo.Collection
	CollUser *mgo.Collection
}

func (cr *CheckInResource) Init(session *mgo.Session) {
	cr.Mongo = session
	cr.CollCI = cr.Mongo.DB("tidy").C("checkin")
	cr.CollUser = cr.Mongo.DB("tidy").C("user")
}

// CheckIn do checkin task for special user id
// Method: POST
func (cr *CheckInResource) CheckIn(c *gin.Context) {
	now := time.Now()
	content := c.PostForm("content")
	uidString := c.PostForm("uid")
	log.Print("Checkin user_id: " + uidString)
	uid := bson.ObjectIdHex(uidString)
	err := cr.CollCI.Insert(&mod.CheckIn{
		Id_:         bson.NewObjectId(),
		UserId:      uid,
		Content:     content,
		CreateAt:    now,
		CreateDay:   now.Day(),
		CreateMonth: int(now.Month()),
		CreateYear:  now.Year(),
		CreateHour:  now.Hour(),
		CreateMin:   now.Minute(),
		CreateSec:   now.Second(),
		Timestamp:   now.Unix(),
		Images:      []string{"abc.png", "xyz.png"},
	})
	if err != nil {
		panic(err)
	}
	err = cr.CollUser.Update(
		bson.M{
			"_id": uid,
		},
		bson.M{
			"$inc": bson.M{
				"continuous": 1,
			},
		})

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, now.Unix())
}

// ListCheckIn return all checkin records
// Method: GET
func (cr *CheckInResource) ListCheckIn(c *gin.Context) {
	col := cr.Mongo.DB("tidy").C("checkin")
	uid := bson.ObjectIdHex(c.DefaultQuery("uid", ""))
	//objectId := bson.ObjectIdHex(id)
	//log.Print("user_id: " + uid)
	var ci []mod.CheckIn
	col.Find(bson.M{"user_id": uid}).All(&ci)
	//col.Find(nil).All(&ci)
	//log.Printf("%s", ci)
	c.JSON(http.StatusOK, ci)
}
