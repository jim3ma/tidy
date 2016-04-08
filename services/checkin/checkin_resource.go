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
	content := c.PostForm("content")
	err := col.Insert(&mod.CheckIn{
		Id_:       bson.NewObjectId(),
		UserId:    bson.NewObjectId(),
		Content:   	content,
		CreateAt:  	now,
		CreateDay: 	now.Day(),
		CreateMonth: 	int(now.Month()),
		CreateYear: 	now.Year(),
		Images:		[]string{"abc.png", "xyz.png"},
	})

	if err != nil {
		panic(err)
	}
	
	c.JSON(200, now.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
}

func (cr *CheckInResource) QueryMouth(c *gin.Context) {
	id, err := strconv.Atoi("3")
	if err != nil {
		log.Print(err)
	}
	log.Print(id)
	c.JSON(200, id)
}
