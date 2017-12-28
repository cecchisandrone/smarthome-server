package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Raspsonar struct {
	ConfigurationService  *service.Configuration                `inject:""`
	RaspsonarService      *service.Raspsonar                    `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (r Raspsonar) InitRoutes() {

	profile := r.Router.Group("/api/v1/configurations/:id/raspsonar").Use(r.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", r.getMeasurements)
	profile.PUT("/relay", r.toggleRelay)
	profile.GET("/relay", r.getRelayStatus)
}

func (r Raspsonar) getMeasurements(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	scheduledMeasurements := ctx.DefaultQuery("scheduled", "false")

	configuration := r.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	if scheduledMeasurements == "false" {
		timestamp, value, err := r.RaspsonarService.GetLast(*configuration)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "value": value})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		}
	} else {
		measurements := r.RaspsonarService.GetScheduledMeasurements()
		ctx.JSON(http.StatusOK, &measurements)
	}
}

func (r Raspsonar) toggleRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	relayStatus := ctx.DefaultQuery("relayStatus", "1")
	status, err := strconv.Atoi(relayStatus)
	if err == nil {
		configuration := r.checkConfiguration(configurationID, ctx)
		if configuration == nil {
			return
		}
		err := r.RaspsonarService.ToggleRelay(*configuration, status)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"relayStatus": r.RaspsonarService.RelayStatus, "activationTime": r.RaspsonarService.RelayActivationTimestamp})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		}
	}
}

func (r Raspsonar) getRelayStatus(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	configuration := r.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"relayStatus": r.RaspsonarService.RelayStatus, "activationTime": r.RaspsonarService.RelayActivationTimestamp})
}

func (r Raspsonar) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := r.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
