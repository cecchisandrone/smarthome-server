package service

import (
	"net/http"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
)

type Profile struct {
	Db     *gorm.DB    `inject:""`
	Router *gin.Engine `inject:""`
}

func (p Profile) InitRoutes() {

	profile := p.Router.Group("/api/v1/profiles")

	profile.GET("/", p.getProfiles)
	profile.GET("/:id", p.getProfile)
	profile.POST("/", p.createProfile)
}

func (p Profile) getProfiles(c *gin.Context) {

	var profiles []model.Profile
	p.Db.Find(&profiles)
	c.JSON(http.StatusOK, profiles)
}

func (p Profile) createProfile(c *gin.Context) {

	var profile model.Profile
	if err := c.ShouldBindWith(&profile, binding.JSON); err == nil {
		p.Db.Save(&profile)
		c.JSON(http.StatusCreated, profile)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (p Profile) getProfile(c *gin.Context) {

	var profile model.Profile
	profileID := c.Param("id")
	p.Db.First(&profile, profileID)
	if profile.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No profile found!"})
		return
	}
	c.JSON(http.StatusOK, profile)
}
