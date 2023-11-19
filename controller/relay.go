package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
)

type Relay struct {
	ConfigurationService  *service.Configuration                `inject:""`
	RelayService          *service.Relay                        `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (w Relay) InitRoutes() {

	profile := w.Router.Group("/api/v1/configurations/:id/relays").Use(w.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.PUT("/:relayId/relay", w.toggleRelay)
	profile.GET("/:relayId/relay", w.getRelayStatus)
	profile.GET("/", w.getRelays)
	profile.POST("/", w.createRelay)
	profile.GET("/:relayId", w.getRelay)
	profile.PUT("/:relayId", w.updateRelay)
	profile.DELETE("/:relayId", w.deleteRelay)
}

func (w Relay) getRelays(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	ctx.JSON(http.StatusOK, w.RelayService.GetRelays(configurationID))
}

func (w Relay) createRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var relay model.Relay
	if err := ctx.ShouldBindWith(&relay, binding.JSON); err == nil {
		w.RelayService.CreateOrUpdateRelay(configurationID, &relay)
		ctx.JSON(http.StatusCreated, relay)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (w Relay) getRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	relayID := ctx.Param("relayId")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	relay, err := w.RelayService.GetRelay(relayID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, relay)
}

func (w Relay) updateRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var relay model.Relay

	if err := ctx.ShouldBindWith(&relay, binding.JSON); err == nil {
		w.RelayService.CreateOrUpdateRelay(configurationID, &relay)
		ctx.JSON(http.StatusAccepted, relay)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (w Relay) deleteRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	relayID := ctx.Param("relayId")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	err := w.RelayService.DeleteRelay(relayID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Deleted")
}

func (w Relay) getRelayStatus(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	relayId := ctx.Param("relayId")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	relay, err := w.RelayService.GetRelay(relayId)

	if err == nil {

		resp, err := w.RelayService.GetRelayStatus(relay)

		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"status": resp})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func (w Relay) toggleRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	relayId := ctx.Param("relayId")

	var body map[int]bool

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	manuallyActivated := ctx.DefaultQuery("manuallyActivated", "false")
	manuallyActivatedBool, err := strconv.ParseBool(manuallyActivated)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unable to parse request: " + err.Error()})
	} else {
		if err := ctx.ShouldBindWith(&body, binding.JSON); err == nil {
			relay, err := w.RelayService.GetRelay(relayId)

			if err == nil {
				err = w.RelayService.ToggleRelay(relay, body, manuallyActivatedBool)
				if err == nil {
					ctx.JSON(http.StatusOK, body)
				} else {
					ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": http.StatusNotFound, "message": err.Error()})
				}
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Unable to parse request: " + err.Error()})
		}
	}
}

func (w Relay) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := w.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
