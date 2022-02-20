package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type PowerMeter struct {
	ConfigurationService  *service.Configuration                `inject:""`
	PowerMeterService     *service.PowerMeter                   `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (p PowerMeter) InitRoutes() {

	profile := p.Router.Group("/api/v1/configurations/:id/power").Use(p.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", p.getMeasurements)
}

func (p PowerMeter) getMeasurements(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	scheduledMeasurements := ctx.DefaultQuery("scheduled", "false")

	configuration := p.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	if scheduledMeasurements == "false" {
		timestamp, value, err := p.PowerMeterService.GetLast(*configuration)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "value": value})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		}
	} else {
		measurements := p.PowerMeterService.GetScheduledMeasurements()
		ctx.JSON(http.StatusOK, &measurements)
	}
}

func (p PowerMeter) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := p.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
