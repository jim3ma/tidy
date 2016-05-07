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

	"github.com/spf13/viper"
)

type UserResource struct {
	Mongo    *mgo.Session
	CollUser *mgo.Collection
}

func (ur *UserResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	ur.Mongo = session
	ur.CollUser = ur.Mongo.DB(db).C("user")
}

// NewUser add a user into mongo/tidy/user
// return current timestamp if success
func (ur *UserResource) NewUser(c *gin.Context) {
	now := time.Now()
	//col := ur.Mongo.DB("tidy").C("user")
	//content := c.PostForm("content")
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	if username == "" || password == "" || email == "" {
		c.JSON(http.StatusBadRequest, "Invalid parameter")
	}
	log.Print("New username:" + username)
	log.Print("New password:" + password)
	log.Print("New email:" + email)
	user := &mod.User{
		ID:         bson.NewObjectId(),
		UserName:   username,
		Password:   util.Md5Sum(password),
		EMail:      email,
		CreateAt:   now,
		Timestamp:  now.Unix(),
		Portrait:   "avantar.png",
		Continuous: 0,
		//LastCheckIn:  ,
	}
	err := ur.CollUser.Insert(user)

	if err != nil {
		panic(err)
	}
	ur.newTokenAndRet(c, user)
}

type AuthReponse struct {
	AuthToken string   `json:"auth_token"`
	UserInfo  mod.User `json:"user_info"`
}

func (ur *UserResource) AuthWithPassword(c *gin.Context) {
	username := c.DefaultQuery("username", "")
	password := c.DefaultQuery("password", "")
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, "Invalid username or password")
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
	ur.newTokenAndRet(c, user)
}

func (ur *UserResource) newTokenAndRet(c *gin.Context, user *mod.User) {
	tokenString, err := util.NewToken(
		map[string]string{
			"uid":       user.ID.Hex(),
			"user_name": user.UserName,
		})
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK,
		AuthReponse{
			AuthToken: tokenString,
			UserInfo:  *user,
		})
}

func (ur *UserResource) QueryInfo(c *gin.Context) {
	uidString := c.DefaultQuery("uid", "")
	if uidString == "" {
		c.JSON(http.StatusBadRequest, "Empty User ID")
		return
	}
	uid := bson.ObjectIdHex(uidString)
	user := new(mod.User)
	err := ur.CollUser.Find(
		bson.M{
			"_id": uid,
		}).One(user)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, struct {
		UserInfo interface{} `json:"user_info"`
	}{
		UserInfo: user,
	})
}
