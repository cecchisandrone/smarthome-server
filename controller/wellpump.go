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

type WellPump struct {
	ConfigurationService  *service.Configuration                `inject:""`
	WellPumpService       *service.WellPump                     `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (w WellPump) InitRoutes() {

	profile := w.Router.Group("/api/v1/configurations/:id/well-pumps").Use(w.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.PUT("/:wellPumpId/relay", w.toggleRelay)
	profile.GET("/:wellPumpId/relay", w.getRelay)
	profile.GET("/", w.getWellPumps)
	profile.POST("/", w.createWellPump)
	profile.GET("/:wellPumpId", w.getWellPump)
	profile.PUT("/:wellPumpId", w.updateWellPump)
	profile.DELETE("/:wellPumpId", w.deleteWellPump)
}

func (w WellPump) getWellPumps(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	ctx.JSON(http.StatusOK, w.WellPumpService.GetWellPumps(configurationID))
}

func (w WellPump) createWellPump(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var wellPump model.WellPump
	if err := ctx.ShouldBindWith(&wellPump, binding.JSON); err == nil {
		w.WellPumpService.CreateOrUpdateWellPump(configurationID, &wellPump)
		ctx.JSON(http.StatusCreated, wellPump)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (w WellPump) getWellPump(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	wellPumpID := ctx.Param("wellPumpId")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	wellPump, err := w.WellPumpService.GetWellPump(wellPumpID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, wellPump)
}

func (w WellPump) updateWellPump(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var wellPump model.WellPump

	if err := ctx.ShouldBindWith(&wellPump, binding.JSON); err == nil {
		w.WellPumpService.CreateOrUpdateWellPump(configurationID, &wellPump)
		ctx.JSON(http.StatusAccepted, wellPump)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (w WellPump) deleteWellPump(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	wellPumpID := ctx.Param("wellPumpId")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	err := w.WellPumpService.DeleteWellPump(wellPumpID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Deleted")
}

func (w WellPump) getRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	wellPumpId := ctx.Param("wellPumpId")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	wellPump, err := w.WellPumpService.GetWellPump(wellPumpId)

	if err == nil {

		statusInt, err := w.WellPumpService.GetRelay(wellPump)

		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"status": statusInt})
		} else {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func (w WellPump) toggleRelay(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	wellPumpId := ctx.Param("wellPumpId")

	if configuration := w.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	status := ctx.DefaultQuery("status", "0")
	statusInt, err := strconv.Atoi(status)

	if err == nil {

		wellPump, err := w.WellPumpService.GetWellPump(wellPumpId)

		if err == nil {

			err = w.WellPumpService.ToggleRelay(wellPump, statusInt)

			if err == nil {
				ctx.JSON(http.StatusOK, gin.H{"status": statusInt})
			} else {
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": http.StatusNotFound, "message": err.Error()})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
}

func (w WellPump) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := w.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
