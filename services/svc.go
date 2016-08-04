package services

import (
	//"log"

	"time"

	"gopkg.in/mgo.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/jim3mar/basicmgo/mongo"
	"github.com/jim3mar/endless"
	jsonp "github.com/jim3mar/gin-jsonp"
	cr "github.com/jim3mar/tidy/services/checkin"
	"github.com/jim3mar/tidy/services/oauth2"
	sr "github.com/jim3mar/tidy/services/system"
	ur "github.com/jim3mar/tidy/services/user"
	util "github.com/jim3mar/tidy/utilities"
	"github.com/jim3mar/tidy/utilities/cache"
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

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	//log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
	//log.SetLevel(log.InfoLevel)
}

func (s *Service) getMgoSession(cfg Config) (*mgo.Session, error) {
	//if bs, err := json.MarshalIndent(cfg, "", "    "); err != nil {
	//	panic(err)
	//} else {
	//	log.Info("Current configuration:" + string(bs))
	//}

	mgoSession, err := mongo.CopyMonotonicSession()
	if err != nil {
		log.Fatalf("CreateMongoSession: %s\n", err)
		return nil, err
	}
	return mgoSession, nil
}

func (s *Service) Run(cfg Config) error {
	go cache.InitCacheConfig()

	util.InitMailConfig()
	mgoSession, err := s.getMgoSession(cfg)

	if err != nil {
		return err
	}
	defer mgoSession.Close()

	svcSR := &sr.SysResource{}
	svcSR.Init(mgoSession)

	svcUR := &ur.UserResource{}
	svcUR.Init(mgoSession)
	svcUR.SystemResource = svcSR

	svcCR := &cr.CheckInResource{}
	svcCR.Init(mgoSession)
	svcCR.UserResource = svcUR

	svcWR := &oauth2.WeChatResource{}
	svcWR.Init(mgoSession)
	svcWR.UserResource = svcUR

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
	router.Use(util.DebugHandler())

	v1 := router.Group("/v1")
	{
		/////////////////////////////////////////
		// checkin api
		// need token
		ci := v1.Group("/checkin")
		ci.Use(util.JWTHandler())

		ci.POST("", svcCR.CheckIn)
		ci.PUT("", svcCR.EditCheckIn)
		ci.GET("", svcCR.ListCheckIn)
		ci.DELETE("", svcCR.DeleteCheckIn)

		ci.PUT("/public", svcCR.MarkCIPublic)
		ci.PUT("/private", svcCR.MarkCIPrivate)

		ci.POST("/favor", svcCR.FavorCheckIn)
		ci.DELETE("/favor", svcCR.UnFavorCheckIn)

		ci.POST("/thumb", svcCR.ThumbCheckIn)
		ci.DELETE("/thumb", svcCR.UnThumbCheckIn)

		ci.POST("/comment", svcCR.CommentCheckIn)
		ci.DELETE("/comment", svcCR.UnCommentCheckIn)
		//ci.POST("/uploadimg", svcCR.UploadImg)

		/////////////////////////////////////////
		o := v1.Group("/oauth2")
		o.GET("/wechat", svcWR.ExchangeToken)
		o.POST("/wechat", svcWR.ExchangeToken)
		o.GET("/wechat_url", svcWR.CreateAuthURL)

		/////////////////////////////////////////
		// user api: register and login
		user := v1.Group("/user")
		user.POST("/uploadimg", svcCR.UploadImg)
		user.POST("/register", svcUR.RegisterUser)
		user.POST("/reset_password", svcUR.ResetPassword)
		user.GET("/query", svcUR.RegisterQuery)
		user.GET("/login", svcUR.AuthWithPassword)

		/////////////////////////////////////////
		user.POST("/feedback", svcSR.CreateFeedback)

		/////////////////////////////////////////
		// user infomation
		// need token
		userInfo := user.Group("/info")
		userInfo.Use(util.JWTHandler())
		userInfo.GET("", svcUR.QueryUserInfo)

		updateSetting := user.Group("/update_setting")
		updateSetting.Use(util.JWTHandler())
		updateSetting.POST("", svcUR.UpdateSetting)
		updateSetting.POST("/portrait", svcUR.UpdatePortrait)
		updateSetting.POST("/password", svcUR.UpdatePassword)

		/////////////////////////////////////////
		// static files
		v1.Static("/static/images", "./tmp")
		//v1.Static("/static", ".")
	}

	//router.Run(cfg.ServiceHost)
	endless.ListenAndServe(cfg.ServiceHost, router)

	return nil
}
