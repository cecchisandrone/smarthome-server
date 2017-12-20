package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type Notification struct {
	ConfigurationService  *service.Configuration                `inject:""`
	NotificationService   *service.Notification                 `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (n Notification) InitRoutes() {

	profile := n.Router.Group("/api/v1/configurations/:id/notification").Use(n.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.POST("/slack/test", n.test)
}

func (n Notification) test(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	configuration := n.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	err := n.NotificationService.SendSlackMessage("alarm", "This is a test message")

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	} else {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": http.StatusServiceUnavailable, "message": err.Error()})
	}
}

func (n Notification) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := n.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
