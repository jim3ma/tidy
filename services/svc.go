package services

import (
	//"encoding/json"
)

type Config struct {
        ServiceHost     string  `json:"service_host"`
        MongoDBHosts    string  `json:"mongo_hosts"`
        MongoAuthUser   string  `json:"mongo_user"`
        MongoAuthPass   string  `json:"mongo_passwd"`
        MongoAuthDB     string  `json:"mongo_database"`
}
