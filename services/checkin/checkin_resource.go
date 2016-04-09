package checkin

import (
	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/checkin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	//"strconv"
	"time"
)

type CheckInResource struct {
	mongo *mgo.Session
}

func (cr *CheckInResource) CheckIn(c *gin.Context) {
	now := time.Now()
	col := cr.mongo.DB("tidy").C("checkin")
	content := c.PostForm("content")
	err := col.Insert(&mod.CheckIn{
		Id_:         bson.NewObjectId(),
		UserId:      bson.NewObjectId(),
		Content:     content,
		CreateAt:    now,
		CreateDay:   now.Day(),
		CreateMonth: int(now.Month()),
		CreateYear:  now.Year(),
		CreateHour:  now.Hour(),
		CreateMin:   now.Minute(),
		CreateSec:   now.Second(),
		Timestamp:	 now.Unix(),
		Images:      []string{"abc.png", "xyz.png"},
	})

	if err != nil {
		panic(err)
	}

	c.JSON(200, now.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
}

func (cr *CheckInResource) ListCheckIn(c *gin.Context) {
	col := cr.mongo.DB("tidy").C("checkin")
	//id := c.DefaultQuery("id", "")
	//objectId := bson.ObjectIdHex(id)
	var ci []mod.CheckIn
	//c.Find(bson.M{"_id": objectId}).One(&ci)
	col.Find(nil).All(&ci)
	log.Printf("%s", ci)
	c.JSON(200, ci)
}
