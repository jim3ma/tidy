package oauth2

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"

        mu "github.com/jim3mar/tidy/models/user"
        mwu "github.com/jim3mar/tidy/models/wechat"
	mpo "github.com/chanxuehong/wechat.v2/mp/oauth2"
	o "github.com/chanxuehong/wechat.v2/oauth2"
	svcuser "github.com/jim3mar/tidy/services/user"

	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
)

type WeChatResource struct {
	Mongo        *mgo.Session
	CollWeChat   *mgo.Collection
	UserResource *svcuser.UserResource

	AppId       string
	AppSecret   string
	RedirectURI string
	Scope       string
	Endpoint    o.Endpoint
}

func (w *WeChatResource) Init(session *mgo.Session) {
	db := viper.GetString("mongo.db")
	w.Mongo = session
	w.CollWeChat = w.Mongo.DB(db).C("wechat")

	w.AppId = "AppId"
	w.AppSecret = "AppSecret"
	w.RedirectURI = "http://api.ctidy.com/oauth/wechat"
	w.Scope = "snsapi_userinfo"
	w.Endpoint = mpo.NewEndpoint(w.AppId, w.AppSecret)
}

func (w *WeChatResource) CreateAuthURL(c *gin.Context) {
	state := "wechat"
	authUrl := mpo.AuthCodeURL(w.AppId, w.RedirectURI, w.Scope, state)
	log.Println("AuthCodeURL:", authUrl)
	c.JSON(http.StatusOK, gin.H{"url": authUrl})
}

func (w *WeChatResource) ExchangeToken(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, "Empth code!")
	}
	client := o.Client{
		Endpoint: w.Endpoint,
	}

	token, err := client.ExchangeToken(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Exchage token failed, code error!")
		return
	}

	userinfo, err := mpo.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		c.JSON(http.StatusNotFound, "Can not get userinfo!")
		return
	}
	ex, chkerr := w.CheckUser(userinfo)
	if chkerr != nil {
		panic(chkerr)	
	}
	
	if ex == false {
		w.CreateUser(userinfo)
	}

	c.JSON(http.StatusOK, userinfo)
}

const (
	WcUserExist = iota
	WcUserNotExist
	TiUserExist
	TiUserNotExist
)

func (w *WeChatResource) CheckUser(rawUser *mpo.UserInfo) (bool, error) {
	var wcUser mwu.WeChatUserInfo
	err := w.CollWeChat.Find(bson.M{"openid": rawUser.OpenId}).One(&wcUser)
	if err != nil {
		return false, err
	}
	if wcUser.OpenId != rawUser.OpenId {
		return false, nil
	}
	return true, nil
}

func (w *WeChatResource) CreateUser(rawUser *mpo.UserInfo) *mu.User {
	b, err := json.Marshal(rawUser)
	if err != nil {
                panic(err)
	}

	var wcUser mwu.WeChatUserInfo
	json.Unmarshal(b, &wcUser)

	wcUser.Id_ = bson.NewObjectId()
	uid := bson.NewObjectId()

	err = w.CollWeChat.Insert(&wcUser)
        if err != nil {
                panic(err)
        }

	now := time.Now()
	user := &mu.User{
                ID:         uid,
                UserName:   wcUser.Nickname,
                Password:   "",
                EMail:      "",
                CreateAt:   now,
                Timestamp:  now.Unix(),
                Portrait:   wcUser.HeadImageURL,
                Continuous: 0,
                //LastCheckIn:  ,
                Setting: mu.Setting{
                                IMGUploadJS: "canvas.js",
                        },
        }
	
	w.UserResource.CreateUser(user)
	return user
}
