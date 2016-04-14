package services

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/jim3mar/basicmgo/mongo"
	jsonp "github.com/jim3mar/gin-jsonp"
	cr "github.com/jim3mar/tidy/services/checkin"
	ur "github.com/jim3mar/tidy/services/user"
	"gopkg.in/mgo.v2"
	"log"
	//"encoding/json"
	//"time"
)

type Config struct {
	ServiceHost   string `json:"service_host"`
	MongoDBHosts  string `json:"mongo_hosts"`
	MongoAuthUser string `json:"mongo_user"`
	MongoAuthPass string `json:"mongo_passwd"`
	MongoAuthDB   string `json:"mongo_database"`
}

type Response struct {
	Status int `json:'status'`
	RedirectTo string `json:redirect_to`
}


type Service struct {
	mgoSession *mgo.Session
}

func (s *Service) getMgoSession(cfg Config) (*mgo.Session, error) {
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

func (s *Service) Run(cfg Config) error {
	mgoSession, err := s.getMgoSession(cfg)

	if err != nil {
		return err
	}
	defer mgoSession.Close()

	svc_cr := &cr.CheckInResource{
		Mongo: mgoSession,
	}

	svc_ur := &ur.UserResource{
		Mongo: mgoSession,
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(jsonp.Handler())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	{
		v1.POST("/checkin", svc_cr.CheckIn)
		v1.GET("/checkin", svc_cr.ListCheckIn)

		v1.GET("/user", svc_ur.NewUser)
	}

	//router.Run(cfg.ServiceHost)
	endless.ListenAndServe(cfg.ServiceHost, router)

	return nil
}
