// Copyright © 2016 Jim Mar
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	//"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jim3ma/basicmgo/mongo"
	"github.com/jim3ma/tidy/services"
	"github.com/jim3ma/tidy/utilities"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start tidy",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Tidy Serve Running...")
		Main()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	// viper.SetEnvPrefix("TIDY")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func updateConfig(config *services.Config) {
	config.ServiceHost = fmt.Sprintf("%s:%s",
		viper.GetString("host"), viper.GetString("port"))
	config.MongoDBHosts = fmt.Sprintf("%s:%s",
		viper.GetString("mongo.host"), viper.GetString("mongo.port"))
	config.MongoAuthUser = viper.GetString("mongo.username")
	config.MongoAuthPass = viper.GetString("mongo.password")
	config.MongoAuthDB = viper.GetString("mongo.db")
	log.Infof("current config: %+v", *config)
}

func updateKeys() {
	if viper.GetString("jwt.pubkey") != "" {
		utilities.UpdatePubKey(viper.GetString("jwt.pubkey"))
	}
	if viper.GetString("jwt.prikey") != "" {
		utilities.UpdatePriKey(viper.GetString("jwt.prikey"))
	}
}

// Main setup mongo session and launch Tidy service
func Main() {
	//cfg, _ := getConfig()
	cfg := services.Config{}
	updateConfig(&cfg)
	updateKeys()

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
	log.Infof("\nCatched Signal: %v\r\n", <-ch)
	log.Infof("Graceful Shutdown.")
}
