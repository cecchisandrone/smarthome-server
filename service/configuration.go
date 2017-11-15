package service

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
)

type Configuration struct {
	Db     *gorm.DB    `inject:""`
	Router *gin.Engine `inject:""`
}

func (c Configuration) InitRoutes() {

	configuration := c.Router.Group("/api/v1/configurations")

	configuration.GET("/", c.getConfigurations)
	configuration.GET("/:id", c.getConfiguration)
	configuration.POST("/", c.createConfiguration)
	configuration.DELETE("/:id", c.deleteConfiguration)
}

func (c *Configuration) getConfigurations(ctx *gin.Context) {

	var configurations []model.Configuration
	c.Db.Find(&configurations)
	ctx.JSON(http.StatusOK, configurations)
}

func (c *Configuration) createConfiguration(ctx *gin.Context) {

	var configuration model.Configuration
	if err := ctx.ShouldBindWith(&configuration, binding.JSON); err == nil {
		c.Db.Save(&configuration)
		ctx.JSON(http.StatusCreated, configuration)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (c *Configuration) getConfiguration(ctx *gin.Context) {

	var configuration model.Configuration
	configurationID := ctx.Param("id")
	c.Db.Preload("Profile").First(&configuration, configurationID)
	if configuration.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No configuration found!"})
		return
	}
	ctx.JSON(http.StatusOK, configuration)
}

func (c *Configuration) deleteConfiguration(ctx *gin.Context) {

	var configuration model.Configuration
	configurationID := ctx.Param("id")
	c.Db.First(&configuration, configurationID)
	if configuration.ID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No configuration found!"})
		return
	}
	c.Db.Unscoped().Delete(&configuration)
	ctx.JSON(http.StatusOK, "Deleted")
}
