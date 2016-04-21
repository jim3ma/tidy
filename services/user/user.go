package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/user"
	util "github.com/jim3mar/tidy/utilities"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"encoding/json"
	"log"
	//"strconv"
	"time"
)

type UserResource struct {
	Mongo    *mgo.Session
	CollUser *mgo.Collection
}

func (ur *UserResource) Init(session *mgo.Session) {
	ur.Mongo = session
	ur.CollUser = ur.Mongo.DB("tidy").C("user")
}

func (ur *UserResource) NewUser(c *gin.Context) {
	now := time.Now()
	//col := ur.Mongo.DB("tidy").C("user")
	//content := c.PostForm("content")
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, "Invalid parameter")
	}
	log.Print(username)
	log.Print(password)
	log.Print(email)
	err := ur.CollUser.Insert(&mod.User{
		Id_:        bson.NewObjectId(),
		UserName:   username,
		Password:   util.Md5Sum(password),
		EMail:      email,
		CreateAt:   now,
		Timestamp:  now.Unix(),
		Portrait:   "avantar.png",
		Continuous: 0,
		//LastCheckIn:  ,
	})

	if err != nil {
		panic(err)
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, now.Unix())
}

type AuthReponse struct {
	AuthToken string   `json:"auth_token"`
	UserInfo  mod.User `json:"user_info"`
}

func (ur *UserResource) AuthWithPassword(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	password := c.DefaultQuery("password", "")
	if username == "" || password == "" {
		c.JSON(http.StatusForbidden, "invalid username or password")
		return
	}
	password = util.Md5Sum(password)
	user := new(mod.User)
	err := ur.CollUser.Find(
		bson.M{
			"username": username,
			"password": password,
		}).One(user)
	if err != nil {
		c.JSON(http.StatusForbidden, err)
		return
	}
	c.JSON(http.StatusOK,
		AuthReponse{
			AuthToken: "570fb03a55cbf50efc93a728",
			UserInfo:  *user,
		})
	return
}
