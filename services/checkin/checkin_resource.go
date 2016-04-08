package checkin

import (
	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/checkin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"time"
)

type CheckInResource struct {
	mongo *mgo.Session
}

func (cr *CheckInResource) CheckIn(c *gin.Context) {
	now := time.Now()
	col := cr.mongo.DB("tidy").C("tidy")
	err := col.Insert(&mod.CheckIn{
		Id_:       bson.NewObjectId(),
		UserId:    bson.NewObjectId(),
		Content:   	"tidy-checkin-content",
		CreateAt:  	now,
		CreateDay: 	now.Day(),
		CreateMonth: 	int(now.Month()),
		CreateYear: 	now.Year(),
		Images:		[]string{"abc.png", "xyz.png"},
	})

	if err != nil {
		panic(err)
	}
	
	c.JSON(200, now.Format("2016-04-08 00:00:00"))
}

func (cr *CheckInResource) QueryMouth(c *gin.Context) {
	id, err := strconv.Atoi("3")
	if err != nil {
		log.Print(err)
	}
	log.Print(id)
	c.JSON(200, id)
}
