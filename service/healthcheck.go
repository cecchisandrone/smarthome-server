package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct {
	Router *gin.Engine `inject:""`
}

func (h HealthCheck) InitRoutes() {
	fmt.Print("assdf")
	h.Router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}
