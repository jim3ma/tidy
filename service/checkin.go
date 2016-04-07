package service

import (
	"gopkg.in/mgo.v2"
	"github.com/gin-gonic/gin"
	"log"
	"encoding/json"
	"time"
)

type Config struct {
	ServiceHost	string	`json:"service_host"`
	MongoDBHosts	string	`json:"mongo_hosts"`
	MongoAuthUser	string  `json:"mongo_user"`
	MongoAuthPass	string  `json:"mongo_passwd"`
	MongoAuthDB	string  `json:"mongo_database"`
}

type CheckInService struct {
}

func (s *CheckInService) getDb(cfg Config) (*mgo.Session, error){
	if bs, err := json.MarshalIndent(cfg, "", "    "); err != nil {
		panic(err)
	} else {
		log.Print("Current configuration:\n" + string(bs))
	}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{cfg.MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: cfg.MongoAuthDB,
		Username: cfg.MongoAuthUser,
		Password: cfg.MongoAuthPass,
	}
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateMongoSession: %s\n", err)
		return nil, err 
	}
	mongoSession.SetMode(mgo.Monotonic, true)
	return mongoSession, nil
}

func (s *CheckInService) Run(cfg Config) error {
	mongoSession, err := s.getDb(cfg)

	if err != nil {
		return err
	}
	defer mongoSession.Close()

	cr := &CheckInResource{mongo: mongoSession}

	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	route.GET("/checkin", cr.CheckIn)

	route.Run(cfg.ServiceHost)

	return nil
}
