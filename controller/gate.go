package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type Gate struct {
	ConfigurationService  *service.Configuration                `inject:""`
	GateService           *service.Gate                         `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (g Gate) InitRoutes() {

	profile := g.Router.Group("/api/v1/configurations/:id/gate").Use(g.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.POST("/open", g.open)
}

func (g Gate) open(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	configuration := g.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	err := g.GateService.Open(*configuration)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
	} else {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	}
}

func (g Gate) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := g.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
