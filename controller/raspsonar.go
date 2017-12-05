package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type Raspsonar struct {
	ConfigurationService  *service.Configuration                `inject:""`
	RaspsonarService      *service.Raspsonar                    `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (t Raspsonar) InitRoutes() {

	profile := t.Router.Group("/api/v1/configurations/:id/raspsonar").Use(t.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", t.getMeasurements)
}

func (t Raspsonar) getMeasurements(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	scheduledMeasurements := ctx.DefaultQuery("scheduled", "false")

	configuration := t.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	if scheduledMeasurements == "false" {
		timestamp, value := t.RaspsonarService.GetLast(*configuration)
		ctx.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "value": value})
	} else {
		measurements := t.RaspsonarService.GetScheduledMeasurements()
		ctx.JSON(http.StatusOK, &measurements)
	}
}

func (t Raspsonar) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := t.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
