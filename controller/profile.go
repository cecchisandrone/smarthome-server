package controller

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/authentication"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Profile struct {
	ProfileService        *service.Profile                      `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (p Profile) InitRoutes() {

	profile := p.Router.Group("/api/v1/profiles").Use(p.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	profile.GET("/", p.getProfiles)
	profile.GET("/:id", p.getProfile)
	profile.POST("/", p.createProfile)
}

func (p Profile) getProfiles(ctx *gin.Context) {

	profiles, err := p.ProfileService.GetProfiles()
	if err == nil {
		ctx.JSON(http.StatusOK, profiles)
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (p Profile) createProfile(ctx *gin.Context) {

	var err error
	var profile model.Profile
	if err = ctx.ShouldBindWith(&profile, binding.JSON); err == nil {
		if err = p.ProfileService.CreateProfile(&profile); err == nil {
			ctx.JSON(http.StatusCreated, profile)
		}
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (p Profile) getProfile(ctx *gin.Context) {

	profileID := ctx.Param("id")
	profile, err := p.ProfileService.GetProfile(profileID)
	if profile == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusNotFound, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, profile)
	}
}
