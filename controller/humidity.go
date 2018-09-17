package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type Humidity struct {
	ConfigurationService  *service.Configuration                `inject:""`
	HumidityService       *service.Humidity                     `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (h Humidity) InitRoutes() {

	profile := h.Router.Group("/api/v1/configurations/:id/humidity").Use(h.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", h.getMeasurements)
}

func (h Humidity) getMeasurements(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	scheduledMeasurements := ctx.DefaultQuery("scheduled", "false")

	configuration := h.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	if scheduledMeasurements == "false" {
		timestamp, value, err := h.HumidityService.GetLast(*configuration)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "value": value})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		}
	} else {
		measurements := h.HumidityService.GetScheduledMeasurements()
		ctx.JSON(http.StatusOK, &measurements)
	}
}

func (h Humidity) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := h.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
