package checkin

import (
	"gopkg.in/mgo.v2"
	"github.com/gin-gonic/gin"
	"github.com/jim3mar/basicmgo/mongo"
	"github.com/jim3mar/tidy/services"
	"log"
	"encoding/json"
	"time"
)

type Service struct {
	mgoSession *mgo.Session
}

func (s *Service) getMgoSession(cfg services.Config) (*mgo.Session, error){
	if bs, err := json.MarshalIndent(cfg, "", "    "); err != nil {
		panic(err)
	} else {
		log.Print("Current configuration:\n" + string(bs))
	}

	mgoconfig := &mongo.MongoConfiguration{
		Hosts:   	cfg.MongoDBHosts,
                Database:	cfg.MongoAuthDB,
                UserName:	cfg.MongoAuthUser,
                Password:	cfg.MongoAuthPass,
                Timeout: 	60 * time.Second,
	}

	if err := mongo.Startup(mgoconfig); err != nil {
		log.Fatalf("MongoSession startup failed: %s\n", err)
		return nil, err 
	}

	mgoSession, err := mongo.CopyMonotonicSession()
	if err != nil {
		log.Fatalf("CreateMongoSession: %s\n", err)
		return nil, err 
	}
	return mgoSession, nil
}

func (s *Service) Run(cfg services.Config) error {
	mgoSession, err := s.getMgoSession(cfg)

	if err != nil {
		return err
	}
	defer mgoSession.Close()

	cr := &CheckInResource{ 
		mongo: mgoSession,
	}

	route := gin.New()
	route.Use(gin.Logger())
	route.Use(gin.Recovery())

	route.GET("/checkin", cr.CheckIn)

	route.Run(cfg.ServiceHost)

	return nil
}
