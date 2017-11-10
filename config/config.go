package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func Init() {
	// Detect environment
	env, envSet := os.LookupEnv("SMARTHOME_ENV")
	if !envSet {
		env = "dev"
	}

	log.Info("SmartHome starting with environment ", env)

	viper.SetConfigName("config/" + env) // no need to include file extension
	viper.AddConfigPath(".")             // set the path of your config file
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Config file not found...", err)
		os.Exit(1)
	}
}
