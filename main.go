package main

import (
	"fmt"
	"os"

	"github.com/cecchisandrone/smarthome-server/config"
	"github.com/cecchisandrone/smarthome-server/controller"
	"github.com/cecchisandrone/smarthome-server/service"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/persistence"
	"github.com/facebookgo/inject"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	controllers := []controller.Controller{&controller.HealthCheck{}, &controller.Profile{}, &controller.Configuration{}, &controller.Camera{}, &controller.Authentication{}}
	services := []service.Service{&service.Profile{}, &service.Configuration{}, &service.Camera{}}

	for _, c := range controllers {
		g.Provide(&inject.Object{Value: c})
	}

	for _, s := range services {
		g.Provide(&inject.Object{Value: s})
	}

	authMiddlewareFactory := &authentication.AuthMiddlewareFactory{}

	g.Provide(&inject.Object{Value: db})
	g.Provide(&inject.Object{Value: router})
	g.Provide(&inject.Object{Value: authMiddlewareFactory})

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

	router.Run()
}
