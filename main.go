package main

import (
	"fmt"
	"os"

	"github.com/cecchisandrone/smarthome-server/config"
	"github.com/cecchisandrone/smarthome-server/service"

	"github.com/cecchisandrone/smarthome-server/persistence"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
)

func main() {

	var g inject.Graph

	// Prepare and inject dependencies
	config.Init()
	db := persistence.Init()
	router := gin.Default()

	services := []service.Service{&service.HealthCheck{}, &service.Profile{}, &service.Configuration{}}

	for _, s := range services {
		g.Provide(&inject.Object{Value: s})
	}

	g.Provide(&inject.Object{Value: db})
	g.Provide(&inject.Object{Value: router})

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Init service routes
	for _, s := range services {
		s.InitRoutes()
	}

	router.Run()
}
