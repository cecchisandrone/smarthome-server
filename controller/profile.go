package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Profile struct {
	ProfileService *service.Profile `inject:""`
	Router         *gin.Engine      `inject:""`
}

func (p Profile) InitRoutes() {

	profile := p.Router.Group("/api/v1/profiles")

	profile.GET("/", p.getProfiles)
	profile.GET("/:id", p.getProfile)
	profile.POST("/", p.createProfile)
}

func (p Profile) getProfiles(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, p.ProfileService.GetProfiles())
}

func (p Profile) createProfile(ctx *gin.Context) {

	var profile model.Profile
	if err := ctx.ShouldBindWith(&profile, binding.JSON); err == nil {
		p.ProfileService.CreateProfile(&profile)
		ctx.JSON(http.StatusCreated, profile)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (p Profile) getProfile(ctx *gin.Context) {

	profileID := ctx.Param("id")
	profile, err := p.ProfileService.GetProfile(profileID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, profile)
}
