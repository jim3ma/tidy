package checkin

import (
	"gopkg.in/mgo.v2"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/jim3mar/basicmgo/mongo"
	"github.com/jim3mar/tidy/services"
	jsonp "github.com/jim3mar/ginjsonp"
	"log"
	//"encoding/json"
	//"time"
)

type Service struct {
	mgoSession *mgo.Session
}

func (s *Service) getMgoSession(cfg services.Config) (*mgo.Session, error){
	//if bs, err := json.MarshalIndent(cfg, "", "    "); err != nil {
	//	panic(err)
	//} else {
	//	log.Print("Current configuration:\n" + string(bs))
	//}

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

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(jsonp.Handler())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	{
        	v1.POST("/checkin", cr.CheckIn)
		v1.GET("/checkin", cr.CheckIn)
    	}

	//router.Run(cfg.ServiceHost)
	endless.ListenAndServe(cfg.ServiceHost, router)

	return nil
}
