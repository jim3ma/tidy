package user

import (
	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/user"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	//"encoding/json"
	//"log"
	//"strconv"
	"time"
)

type UserResource struct {
	Mongo *mgo.Session
}

func (ur *UserResource) NewUser(c *gin.Context) {
	now := time.Now()
	col := ur.Mongo.DB("tidy").C("user")
	//content := c.PostForm("content")
	err := col.Insert(&mod.User{
		Id_:          bson.NewObjectId(),
		UserName:     "tidy",
		Password:     "tidy",
		EMail:        "tidy@tidy.com",
		CreateAt:     now,
		Timestamp:    now.Unix(),
                Portrait:     "avantar.png",
                Continuous:   0,
                //LastCheckIn:  ,
	})

	if err != nil {
		panic(err)
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, now.Unix())
}

