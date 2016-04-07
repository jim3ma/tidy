package main

import (
	"github.com/majinjing3/neaten-checkin/service"
)

func getConfig() (service.Config, error) {
	config := service.Config{}
	config.ServiceHost = "127.0.0.1:8080"
	config.MongoDBHosts = "majin.xyz:37017"
	config.MongoAuthUser = "ci"
	config.MongoAuthPass = "111111"
	config.MongoAuthDB = "test"
	return config, nil
}

func main() {
	cfg, _ := getConfig()
	svc := service.CheckInService{}
	svc.Run(cfg)
}
