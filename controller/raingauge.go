package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type RainGauge struct {
	ConfigurationService  *service.Configuration                `inject:""`
	TemperatureService    *service.RainGauge                    `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (r RainGauge) InitRoutes() {

	profile := r.Router.Group("/api/v1/configurations/:id/rain-gauge").Use(r.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", r.getMeasurements)
}

func (r RainGauge) getMeasurements(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	scheduledMeasurements := ctx.DefaultQuery("scheduled", "false")

	configuration := r.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	if scheduledMeasurements == "false" {
		timestamp, value := r.TemperatureService.GetLast24hTotal()
		ctx.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "value": value})
	} else {
		measurements := r.TemperatureService.GetScheduledMeasurements()
		ctx.JSON(http.StatusOK, &measurements)
	}
}

func (r RainGauge) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := r.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
