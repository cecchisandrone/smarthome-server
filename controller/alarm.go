package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Alarm struct {
	AlarmService          *service.Alarm                        `inject:""`
	ConfigurationService  *service.Configuration                `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (a Alarm) InitRoutes() {

	camera := a.Router.Group("/api/v1/configurations/:id/alarm").Use(a.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	camera.GET("/", a.getAlarmStatus)
	camera.PUT("/", a.toggleAlarm)
}

func (a Alarm) toggleAlarm(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	statusString := ctx.DefaultQuery("status", "1")
	status, err := strconv.Atoi(statusString)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}

	configuration := a.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	cameras, err := a.AlarmService.ToggleAlarm(*configuration, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "cameras": cameras})
}

func (a Alarm) getAlarmStatus(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	configuration := a.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "alarmEnabled": a.AlarmService.AlarmStatus, "cameras": a.AlarmService.Cameras})
}

func (a Alarm) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := a.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
