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
	//"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "tidy",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.tidy.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	setDefault()
	bindEnv()

	viper.SetConfigName(".tidy") // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AddConfigPath(".")     // adding current directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}

func setDefault() {
	viper.SetDefault("host", "0.0.0.0")
	viper.SetDefault("port", "8089")
	viper.SetDefault("debug", "false")
	viper.SetDefault("upload.image", "./tmp/")
	viper.SetDefault("upload.portrait", "./tmp/")
	viper.SetDefault("mongo.host", "127.0.0.1")
	viper.SetDefault("mongo.port", "27017")
	viper.SetDefault("mongo.username", "tidy")
	viper.SetDefault("mongo.password", "111111")
	viper.SetDefault("mongo.db", "tidy")
	viper.SetDefault("user.auth.expire", "120")
	viper.SetDefault("redis.addr", "0.0.0.0:6379")
	viper.SetDefault("redis.passwd", "111111")
	viper.SetDefault("mail.sendname", "no-reply")
}

func bindEnv() {
	viper.BindEnv("host", "TIDY_HOST")
	viper.BindEnv("port", "TIDY_PORT")
	viper.BindEnv("debug", "TIDY_DEBUG")

	// upload folder
	viper.BindEnv("upload.image", "TIDY_UPLOAD_IMAGE")

	// mongo
	viper.BindEnv("mongo.host", "TIDY_MONGODB_HOST")
	viper.BindEnv("mongo.port", "TIDY_MONGODB_PORT")
	viper.BindEnv("mongo.username", "TIDY_MONGODB_USERNAME")
	viper.BindEnv("mongo.password", "TIDY_MONGODB_PASSWORD")
	viper.BindEnv("mongo.db", "TIDY_MONGODB_DATABASE")

	// jwt
	viper.BindEnv("jwt.pubkey", "TIDY_JWT_PUBKEY_PATH")
	viper.BindEnv("jwt.prikey", "TIDY_JWT_PRIKEY_PATH")

	// user
	viper.BindEnv("user.auth.expire", "TIDY_USER_EXPIRE")

	// mail
	viper.BindEnv("mail.host", "TIDY_MAIL_HOST")
	viper.BindEnv("mail.port", "TIDY_MAIL_PORT")
	viper.BindEnv("mail.authaddr", "TIDY_MAIL_AUTH_ADDR")
	viper.BindEnv("mail.authpasswd", "TIDY_MAIL_AUTH_PASSWD")
	viper.BindEnv("mail.sendfrom", "TIDY_MAIL_SEND_FROM")
	viper.BindEnv("mail.sendname", "TIDY_MAIL_SEND_NAME")
	viper.BindEnv("mail.tlsskipverify", "TIDY_MAIL_TLS_SKIP_VERIFY")

	//redis
	viper.BindEnv("redis.addr", "TIDY_REDIS_ADDR")
	viper.BindEnv("redis.passwd", "TIDY_REDIS_PASSWD")
}
