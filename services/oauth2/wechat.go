package oauth2

import (
	//"fmt"
	"encoding/json"
	"log"
	"net/http"
	"time"

	mpo "github.com/jim3mar/wechat.v2/mp/oauth2"
	o "github.com/jim3mar/wechat.v2/oauth2"
	mu "github.com/jim3mar/tidy/models/user"
	mwu "github.com/jim3mar/tidy/models/wechat"
	svcuser "github.com/jim3mar/tidy/services/user"
	util "github.com/jim3mar/tidy/utilities"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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

	w.AppId = "wx2df1e835cf8fef8f"
	w.AppSecret = "c02e50b35b81ac493711d9929defbb58"
	w.RedirectURI = "http://api.ctidy.com/v1/oauth/wechat"
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
	log.Printf("code: %s", code)
	if code == "" {
		code = c.PostForm("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, "Empth code!")
		}
	}
	client := o.Client{
		Endpoint: w.Endpoint,
	}

	token, err := client.ExchangeToken(code)
	log.Printf("token: %s, err: %s", token, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Exchage token failed, code error!")
		return
	}

	userinfo, err := mpo.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	//userinfo, err := mpo.GetUserInfo("OezXcEiiBSKSxW0eoylIeAs5Md6Fpld34iFDYCBQq8sCIPv0MqBa7Z4bjiHxdYKvtNZUkJzwdsAtKOouwmuK-lh7x2wmOrIji_F8b41mCNfTffqX3oBcRUylCYhDFN8s", "oRYmewVi4lTN5dEc_RquC1fqMZ3k", "", nil)
	log.Printf("user info: %s, err:%s", userinfo, err)
	if err != nil {
		c.JSON(http.StatusNotFound, "Can not get userinfo!")
		return
	}
	wcUser, chkerr := w.CheckUser(userinfo)
	if chkerr != nil {
		panic(chkerr)
	}

	var user *mu.User
	if wcUser == nil {
		user = w.CreateUser(userinfo)
	} else {
		user = w.QueryUser(wcUser)
	}
	w.UserResource.CreateTokenR(c, user)
	//c.JSON(http.StatusOK, userinfo)
}

const (
	WcUserExist = iota
	WcUserNotExist
	TiUserExist
	TiUserNotExist
)

func (w *WeChatResource) CheckUser(rawUser *mpo.UserInfo) (*mwu.WeChatUserInfo, error) {
	var wcUser mwu.WeChatUserInfo
	err := w.CollWeChat.Find(bson.M{"openid": rawUser.OpenId}).One(&wcUser)
	log.Printf("check user: %s", err)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	if err != nil && err.Error() == "not found" {
		return nil, nil
	}

	if wcUser.OpenId != rawUser.OpenId {
		return nil, nil
	}
	return &wcUser, nil
}

func (w *WeChatResource) QueryUser(wcUser *mwu.WeChatUserInfo) *mu.User {
	user, err := w.UserResource.QueryUserInfoByID(wcUser.UserId.Hex())
	if err != nil {
		panic(err)
	}
	return user
}

func (w *WeChatResource) CreateUser(rawUser *mpo.UserInfo) *mu.User {
	b, err := json.Marshal(rawUser)
	if err != nil {
		panic(err)
	}

	var wcUser mwu.WeChatUserInfo
	json.Unmarshal(b, &wcUser)
	log.Printf("WcUser: %s", wcUser)

	wcUser.Id_ = bson.NewObjectId()
	uid := bson.NewObjectId()
	wcUser.UserId = uid

	err = w.CollWeChat.Insert(&wcUser)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	username := wcUser.Nickname + "_" + string(util.Krand(8, util.KC_RAND_KIND_ALL))
	user := &mu.User{
		ID:         uid,
		UserName:   username,
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

	log.Printf("TiUser: %s", user)
	w.UserResource.CreateUser(user)
	return user
}
