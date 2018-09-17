package main

import (
	"fmt"
	"os"

	"github.com/cecchisandrone/smarthome-server/config"
	"github.com/cecchisandrone/smarthome-server/controller"
	"github.com/cecchisandrone/smarthome-server/service"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/persistence"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	"github.com/facebookgo/inject"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
	"time"
)

func main() {

	var g inject.Graph

	// Prepare and inject dependencies
	config.Init()
	db := persistence.Init()
	router := gin.Default()
	// CORS config
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Origin", "Content-Length", "Content-Type", "Authorization")
	config.AddAllowMethods("PUT", "DELETE", "GET", "POST")
	router.Use(cors.New(config))

	// Resty config
	// Retries are configured per client
	resty.DefaultClient.SetTimeout(10 * time.Second)

	controllers := []controller.Controller{&controller.HealthCheck{}, &controller.Profile{}, &controller.Configuration{}, &controller.Camera{}, &controller.Authentication{}, &controller.Temperature{}, &controller.Raspsonar{}, &controller.Gate{}, &controller.Notification{}, &controller.Alarm{}, &controller.WellPump{}, &controller.RainGauge{}, &controller.Humidity{}}
	services := []service.Service{&service.Profile{}, &service.Configuration{}, &service.Camera{}, &service.Temperature{}, &service.Raspsonar{}, &service.Gate{}, &service.Notification{}, &service.Alarm{}, &service.WellPump{}, &service.RainGauge{}, &service.Humidity{}}

	for _, c := range controllers {
		g.Provide(&inject.Object{Value: c})
	}

	for _, s := range services {
		g.Provide(&inject.Object{Value: s})
	}

	authMiddlewareFactory := &authentication.AuthMiddlewareFactory{}
	schedulerManager := &scheduler.SchedulerManager{}

	g.Provide(&inject.Object{Value: db})
	g.Provide(&inject.Object{Value: router})
	g.Provide(&inject.Object{Value: authMiddlewareFactory})
	g.Provide(&inject.Object{Value: schedulerManager})

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Init auth middleware
	authMiddlewareFactory.Init()

	// Init controller routes
	for _, c := range controllers {
		c.InitRoutes()
	}

	// Init services
	for _, s := range services {
		s.Init()
	}

	// Start task scheduler
	schedulerManager.Start()

	router.Run()
}
