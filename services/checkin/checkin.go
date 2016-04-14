package checkin

import (
	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/checkin"
	//"github.com/jim3mar/tidy/models/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	//"encoding/json"
	//"log"
	//"strconv"
	"time"
)

type CheckInResource struct {
	Mongo *mgo.Session
	CollCI *mgo.Collection
	CollUser *mgo.Collection
}

func (cr *CheckInResource) Init(session *mgo.Session){
	cr.Mongo = session
	cr.CollCI = cr.Mongo.DB("tidy").C("checkin")
	cr.CollUser = cr.Mongo.DB("tidy").C("user")
}

func (cr *CheckInResource) CheckIn(c *gin.Context) {
	now := time.Now()
	content := c.PostForm("content")
	auth_token := c.PostForm("auth_token")
	// TBD
	auth_token = "570fb03a55cbf50efc93a728"
    uid := bson.ObjectIdHex(auth_token)
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
	//var u = new(user.User)
	//err = cr.CollUser.Find(bson.M{"_id": uid}).One(&u)
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

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, now.Unix())
}

func (cr *CheckInResource) ListCheckIn(c *gin.Context) {
	col := cr.Mongo.DB("tidy").C("checkin")
	//id := c.DefaultQuery("id", "")
	//objectId := bson.ObjectIdHex(id)
	var ci []mod.CheckIn
	//c.Find(bson.M{"_id": objectId}).One(&ci)
	col.Find(nil).All(&ci)
	//log.Printf("%s", ci)
	c.JSON(http.StatusOK, ci)
}
