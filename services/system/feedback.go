package system

import (
	//"errors"
	//"log"
	"net/http"
	//"strconv"
	//"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	ms "github.com/jim3mar/tidy/models/system"
	//util "github.com/jim3mar/tidy/utilities"
	//"github.com/spf13/viper"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (sr *SystemResource) CreateFeedback(c *gin.Context) {
	now := time.Now()
	content := c.PostForm("content")
	username := c.PostForm("user_name")
	log.Infof("username: %s", username)
	log.Infof("content: %s", content)
	if content == "" {
		c.JSON(http.StatusBadRequest, "Empty content")
		return
	}
	fd := &ms.Feedback{
		ID:        bson.NewObjectId(),
		UserName:  username,
		Content:   content,
		Timestamp: now.Unix(),
	}
	err := sr.CollFeedback.Insert(fd)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, "")
}
