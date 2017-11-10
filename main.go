package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

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

	dbInit()
	router := gin.Default()
	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.GET("/", fetchAllTodos)
	}
	router.GET("/health", healthCheck)
	router.Run()
}
