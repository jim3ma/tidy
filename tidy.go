package main

import (
	"github.com/jim3mar/basicmgo/mongo"
	"github.com/jim3mar/tidy/services"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func getConfig() (services.Config, error) {
	config := services.Config{}
	config.ServiceHost = "0.0.0.0:8089"
	config.MongoDBHosts = "127.0.0.1:27017"
	config.MongoAuthUser = "tidy"
	config.MongoAuthPass = "111111"
	config.MongoAuthDB = "tidy"
	return config, nil
}

func main() {
	cfg, _ := getConfig()

	mgocfg := &mongo.MongoConfiguration{
		Hosts:    cfg.MongoDBHosts,
		Database: cfg.MongoAuthDB,
		UserName: cfg.MongoAuthUser,
		Password: cfg.MongoAuthPass,
		Timeout:  60 * time.Second,
	}

	if err := mongo.Startup(mgocfg); err != nil {
		log.Fatalf("MongoSession startup failed: %s\n", err)
		return
	}

	svc := services.Service{}
	go func() {
		svc.Run(cfg)
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("\nCatched Signal: %v\r\n", <-ch)
	log.Printf("Graceful Shutdown.")
}
