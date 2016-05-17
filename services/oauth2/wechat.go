package oauth2

import (
	//"fmt"
	"gopkg.in/mgo.v2"
	"net/http"
	"github.com/gin-gonic/gin"
	mpo "github.com/chanxuehong/wechat.v2/mp/oauth2"
	o "github.com/chanxuehong/wechat.v2/oauth2"
	"log"
)

type WeChatResource struct {
	Mongo       *mgo.Session
	AppId       string
	AppSecret   string
	RedirectURI string
	Scope       string
	Endpoint    o.Endpoint
}

func (w *WeChatResource) Init(session *mgo.Session) {
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

	c.JSON(http.StatusOK, userinfo)
}
