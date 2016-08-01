package cache

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/garyburd/redigo/redis"
	maps "github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

// Pool for redis cache
var Pool *redis.Pool

var redisServer string
var redisPassword string

// InitCacheConfig setup cache configuration
func InitCacheConfig() {
	redisServer = viper.GetString("redis.addr")
	redisPassword = viper.GetString("redis.passwd")
	Pool = newPool(redisServer, redisPassword)
	//log.Infof("current mail config: %+v", config)
}

func dummy() {
	type Person struct {
		Name   string
		Age    int
		Emails []string
		Extra  map[string]string
	}

	// This input can come from anywhere, but typically comes from
	// something like decoding JSON where we're not quite sure of the
	// struct initially.
	input := map[string]interface{}{
		"name":   "Mitchell",
		"age":    91,
		"emails": []string{"one", "two", "three"},
		"extra": map[string]string{
			"twitter": "mitchellh",
		},
	}

	var result Person
	err := maps.Decode(input, &result)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", result)

	///////
	type Server struct {
		Name        string `json:"name,omitempty"`
		ID          int
		Enabled     bool
		users       []string // not exported
		http.Server          // embedded
	}

	server := &Server{
		Name:    "gopher",
		ID:      123456,
		Enabled: true,
	}
	// Convert a struct to a map[string]interface{}
	// => {"Name":"gopher", "ID":123456, "Enabled":true}
	m := structs.Map(server)
	fmt.Printf("%#v", m)

}
