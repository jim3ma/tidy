package service

import (
	"gopkg.in/mgo.v2"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

type CheckInResource struct {
	mongo *mgo.Session
}

func (cr *CheckInResource) CheckIn(c *gin.Context) {
	t := time.Now().Format("Mon Jan 2 15:04:05 +0800 UTC 2006")
	log.Print(t)
	c.JSON(200, t)
}

func (cr *CheckInResource) QueryMouth(c *gin.Context) {
	id, err := strconv.Atoi("3")
	if err != nil {
		log.Print(err)
	}
	log.Print(id)
	c.JSON(200, id)
}
