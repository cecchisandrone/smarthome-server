package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Configuration struct {
	ConfigurationService *service.Configuration `inject:""`
	Router               *gin.Engine            `inject:""`
}

func (c Configuration) InitRoutes() {

	configuration := c.Router.Group("/api/v1/configurations")

	configuration.GET("/", c.getConfigurations)
	configuration.GET("/:id", c.getConfiguration)
	configuration.POST("/", c.createConfiguration)
	configuration.DELETE("/:id", c.deleteConfiguration)
}

func (c Configuration) getConfigurations(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, c.ConfigurationService.GetConfigurations())
}

func (c Configuration) createConfiguration(ctx *gin.Context) {

	var configuration model.Configuration
	if err := ctx.ShouldBindWith(&configuration, binding.JSON); err == nil {
		c.ConfigurationService.CreateConfiguration(&configuration)
		ctx.JSON(http.StatusCreated, configuration)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (c Configuration) getConfiguration(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	configuration, err := c.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, configuration)
}

func (c Configuration) deleteConfiguration(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	if err := c.ConfigurationService.DeleteConfiguration(configurationID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Deleted")
}
