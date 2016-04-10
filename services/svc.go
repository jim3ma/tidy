package services

import (
	//"encoding/json"
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
