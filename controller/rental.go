package controller

import (
	"github.com/cecchisandrone/smarthome-server/dto"
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
)

type Rental struct {
	ConfigurationService  *service.Configuration                `inject:""`
	Router                *gin.Engine                           `inject:""`
	RentalService         *service.Rental                       `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (r Rental) InitRoutes() {

	profile := r.Router.Group("/api/v1/configurations/:id/rental").Use(r.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.POST("/generate-access-link", r.GenerateAccessLink)
}

func (r Rental) GenerateAccessLink(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	configuration := r.checkConfiguration(configurationID, ctx)
	if configuration == nil {
		return
	}

	var booking dto.Booking
	if err := ctx.ShouldBindWith(&booking, binding.JSON); err == nil {
		link, err := r.RentalService.GenerateAccessLink(*configuration, booking)
		if err == nil {
			ctx.JSON(http.StatusOK, gin.H{"link": link})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (r Rental) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := r.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
