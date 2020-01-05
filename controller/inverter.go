package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Inverter struct {
	ConfigurationService  *service.Configuration                `inject:""`
	InverterService       *service.Inverter`inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (i Inverter) InitRoutes() {

	profile := i.Router.Group("/api/v1/configurations/:id/inverters").Use(i.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", i.getInverters)
	profile.POST("/", i.createInverter)
	profile.GET("/:inverterId", i.getInverter)
	profile.PUT("/:inverterId", i.updateInverter)
	profile.DELETE("/:inverterId", i.deleteInverter)
	profile.GET("/:inverterId/metrics", i.getMetrics)
}

func (i Inverter) getInverters(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := i.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	ctx.JSON(http.StatusOK, i.InverterService.GetInverters(configurationID))
}

func (i Inverter) createInverter(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := i.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var Inverter model.Inverter
	if err := ctx.ShouldBindWith(&Inverter, binding.JSON); err == nil {
		i.InverterService.CreateOrUpdateInverter(configurationID, &Inverter)
		ctx.JSON(http.StatusCreated, Inverter)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (i Inverter) getInverter(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	InverterID := ctx.Param("inverterId")

	if configuration := i.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	Inverter, err := i.InverterService.GetInverter(InverterID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Inverter)
}

func (i Inverter) updateInverter(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := i.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var Inverter model.Inverter

	if err := ctx.ShouldBindWith(&Inverter, binding.JSON); err == nil {
		i.InverterService.CreateOrUpdateInverter(configurationID, &Inverter)
		ctx.JSON(http.StatusAccepted, Inverter)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (i Inverter) deleteInverter(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	InverterID := ctx.Param("inverterId")

	if configuration := i.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	err := i.InverterService.DeleteInverter(InverterID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Deleted")
}

func (i Inverter) getMetrics(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	inverterId := ctx.Param("inverterId")

	if configuration := i.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	inverter, err := i.InverterService.GetInverter(inverterId)
	if err == nil {
		value, err := i.InverterService.GetMetrics(inverter)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"metrics": value})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": http.StatusServiceUnavailable, "message": err.Error()})
		}
    } else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func (i Inverter) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := i.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
