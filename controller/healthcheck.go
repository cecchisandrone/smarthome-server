package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheck struct {
	Router *gin.Engine `inject:""`
}

func (h HealthCheck) InitRoutes() {

	h.Router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})
}
