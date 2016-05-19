package user

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	mod "github.com/jim3mar/tidy/models/user"
	util "github.com/jim3mar/tidy/utilities"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

// RegisterUser add a user into mongo/tidy/user
// return current timestamp if success
func (ur *UserResource) RegisterUser(c *gin.Context) {
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

	if ur.IsAccountExist(username, email) {
		c.JSON(http.StatusBadRequest, "username or email account existed!")
		return
	}

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
		Setting: mod.Setting{
			IMGUploadJS: "canvas.js",
		},
	}
	ur.CreateUser(user)
	ur.CreateTokenR(c, user)
}

func (ur *UserResource) CreateUser(user *mod.User) {
	err := ur.CollUser.Insert(user)

	if err != nil {
		panic(err)
	}
}

func (ur *UserResource) IsAccountExist(username string, email string) bool {
	user1, err1 := ur.QueryUserInfoByName(username)
	if err1 != nil || len(user1) != 0 {
		return true
	}

	user2, err2 := ur.QueryUserInfoByEmail(email)
	if err2 != nil || len(user2) != 0 {
		return true
	}
	return false
}

type AuthReponse struct {
	AuthToken string   `json:"auth_token"`
	UserInfo  mod.User `json:"user_info"`
}

func (ur *UserResource) RegisterQuery(c *gin.Context) {
	querytype := c.DefaultQuery("type", "username")
	switch querytype {
	case "username":
		username := c.DefaultQuery("username", "")
		if username != "" {
			user, err := ur.QueryUserInfoByName(username)
			if err != nil || len(user) != 0 {
				c.JSON(http.StatusConflict, "")
				return
			}
		}
	case "email":
		email := c.DefaultQuery("email", "")
		if email != "" {
			user, err := ur.QueryUserInfoByEmail(email)
			if err != nil || len(user) != 0 {
				c.JSON(http.StatusConflict, "")
				return
			}
		}
	default:
		c.JSON(http.StatusConflict, "")
		return
	}
	c.JSON(http.StatusOK, "")
}

func (ur *UserResource) AuthWithPassword(c *gin.Context) {
	account := c.DefaultQuery("account", "")
	password := c.DefaultQuery("password", "")
	if account == "" || password == "" {
		c.JSON(http.StatusBadRequest, "Invalid username, email or password")
		return
	}
	password = util.Md5Sum(password)
	user := new(mod.User)
	var err error
	if strings.Index(account, "@") != -1 {
		err = ur.CollUser.Find(
			bson.M{
				"email":    account,
				"password": password,
			}).One(user)
	} else {
		err = ur.CollUser.Find(
			bson.M{
				"user_name": account,
				"password":  password,
			}).One(user)
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	ur.CreateTokenR(c, user)
}

func (ur *UserResource) CreateToken(user *mod.User) string {
	tokenString, err := util.NewToken(
		map[string]string{
			"uid":       user.ID.Hex(),
			"user_name": user.UserName,
		})
	if err != nil {
		panic(err)
	}

	return tokenString
}

// CreateTokenR create a new token with special user,
// and put the response into c *gin.Context
func (ur *UserResource) CreateTokenR(c *gin.Context, user *mod.User) {
	c.JSON(http.StatusOK,
		AuthReponse{
			AuthToken: ur.CreateToken(user),
			UserInfo:  *user,
		})
}

func (ur *UserResource) QueryUserInfo(c *gin.Context) {
	uidString := c.DefaultQuery("uid", "")
	if uidString == "" {
		c.JSON(http.StatusBadRequest, "Empty User ID")
		return
	}
	user, err := ur.QueryUserInfoByID(uidString)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, struct {
		UserInfo interface{} `json:"user_info"`
	}{
		UserInfo: user,
	})
}

func (ur *UserResource) QueryUserInfoByID(uidString string) (*mod.User, error) {
	if uidString == "" {
		return nil, errors.New("Empty User ID")
	}
	uid := bson.ObjectIdHex(uidString)
	user := new(mod.User)
	err := ur.CollUser.Find(
		bson.M{
			"_id": uid,
		}).One(user)
	return user, err
}

func (ur *UserResource) QueryUserInfoByEmail(email string) ([]mod.User, error) {
	if email == "" {
		return nil, errors.New("Empty email")
	}

	var user []mod.User
	err := ur.CollUser.Find(
		bson.M{
			"email": email,
		}).All(&user)
	return user, err
}

func (ur *UserResource) QueryUserInfoByName(username string) ([]mod.User, error) {
	if username == "" {
		return nil, errors.New("Empty username")
	}

	var user []mod.User
	err := ur.CollUser.Find(
		bson.M{
			"user_name": username,
		}).All(&user)
	return user, err
}

func (ur *UserResource) queryUserHelp(query bson.M, pdata interface{}) error {
	err := ur.CollUser.Find(query).One(pdata)
	return err
}
