package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type Temperature struct {
	ConfigurationService  *service.Configuration                `inject:""`
	TemperatureService    *service.Temperature                  `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (t Temperature) InitRoutes() {

	profile := t.Router.Group("/api/v1/configurations/:id/temperature").Use(t.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", t.getMeasurements)
}

func (t Temperature) getMeasurements(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	scheduledMeasurements := ctx.DefaultQuery("scheduled", "false")

	configuration := t.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	if scheduledMeasurements == "false" {
		timestamp, value, err := t.TemperatureService.GetLast(*configuration)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "value": value})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		}
	} else {
		measurements := t.TemperatureService.GetScheduledMeasurements()
		ctx.JSON(http.StatusOK, &measurements)
	}
}

func (t Temperature) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := t.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
