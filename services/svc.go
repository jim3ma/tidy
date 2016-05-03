package services

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"

	//"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/jim3mar/basicmgo/mongo"
	jsonp "github.com/jim3mar/gin-jsonp"
	cr "github.com/jim3mar/tidy/services/checkin"
	ur "github.com/jim3mar/tidy/services/user"
	"github.com/jim3mar/tidy/utilities"
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
	Status     int    `json:"status"`
	RedirectTo string `json:"redirect_to"`
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

	svcCR := &cr.CheckInResource{}
	svcCR.Init(mgoSession)

	svcUR := &ur.UserResource{}
	svcUR.Init(mgoSession)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	router.Use(jsonp.Handler())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	{
		// checkin api
		// need token
		ci := v1.Group("/checkin")
		ci.Use(utilities.JWTHandler())
		ci.POST("", svcCR.CheckIn)
		//ci.POST("/uploadimg", svcCR.UploadImg)
		ci.GET("", svcCR.ListCheckIn)

		// user api: register and login
		user := v1.Group("/user")
		user.POST("/uploadimg", svcCR.UploadImg)
		user.POST("/register", svcUR.NewUser)
		user.GET("/login", svcUR.AuthWithPassword)

		// user infomation
		// need token
		userInfo := user.Group("/info")
		userInfo.Use(utilities.JWTHandler())
		userInfo.GET("", svcUR.QueryInfo)

		// static files
		v1.Static("/static/images", "./tmp")
		//v1.Static("/static", ".")
	}

	router.Run(cfg.ServiceHost)
	//endless.ListenAndServe(cfg.ServiceHost, router)

	return nil
}
