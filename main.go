package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	dbInit()
	router := gin.Default()
	v1 := router.Group("/api/v1/todos")
	{
		v1.GET("/health", healthCheck)
		v1.POST("/", createTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.GET("/", fetchAllTodos)
	}
	router.Static("/assets", "./assets")
	router.Run()
}
