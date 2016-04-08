package main

import (
	"github.com/jim3mar/tidy/services/checkin"
	"github.com/jim3mar/tidy/services"
	//"github.com/jim3mar/basicmgo/mongo"
)

func getConfig() (services.Config, error) {
	config := services.Config{}
	config.ServiceHost = "10.202.240.252:8089"
	config.MongoDBHosts = "127.0.0.1:27017"
	config.MongoAuthUser = "tidy"
	config.MongoAuthPass = "111111"
	config.MongoAuthDB = "tidy"
	return config, nil
}

func main() {
	cfg, _ := getConfig()
	svc := checkin.Service{}
	svc.Run(cfg)
}
