package main

import (
	"fmt"
	"os"

	"github.com/cecchisandrone/smarthome-server/config"

	"github.com/cecchisandrone/smarthome-server/persistence"
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
)

func main() {

	config.Init()
	db := persistence.Init()

	var g inject.Graph
	s := Service{}

	if err := g.Provide(
		&inject.Object{Value: &s},
		&inject.Object{Value: db},
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	router := gin.Default()
	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", s.createTodo)
		v1.GET("/:id", s.fetchSingleTodo)
		v1.GET("/", s.fetchAllTodos)
	}
	router.GET("/health", s.healthCheck)
	router.Run()
}
