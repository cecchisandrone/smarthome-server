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

type Camera struct {
	CameraService         *service.Camera                       `inject:""`
	ConfigurationService  *service.Configuration                `inject:""`
	Router                *gin.Engine                           `inject:""`
	AuthMiddlewareFactory *authentication.AuthMiddlewareFactory `inject:""`
}

func (c Camera) InitRoutes() {

	camera := c.Router.Group("/api/v1/configurations/:id/cameras").Use(c.AuthMiddlewareFactory.AuthMiddleware.MiddlewareFunc())

	camera.GET("/", c.getCameras)
	camera.POST("/", c.createCamera)
	camera.GET("/:cameraId", c.getCamera)
	camera.PUT("/:cameraId", c.updateCamera)
	camera.PUT("/:cameraId/toggle-alarm", c.toggleCameraAlarm)
	camera.DELETE("/:cameraId", c.deleteCamera)
}

func (c Camera) getCameras(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := c.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	ctx.JSON(http.StatusOK, c.CameraService.GetCameras(configurationID))
}

func (c Camera) createCamera(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := c.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var camera model.Camera
	if err := ctx.ShouldBindWith(&camera, binding.JSON); err == nil {
		c.CameraService.CreateOrUpdateCamera(configurationID, &camera)
		ctx.JSON(http.StatusCreated, camera)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (c Camera) getCamera(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	cameraID := ctx.Param("cameraId")

	if configuration := c.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	camera, err := c.CameraService.GetCamera(cameraID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, camera)
}

func (c Camera) updateCamera(ctx *gin.Context) {

	configurationID := ctx.Param("id")

	if configuration := c.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	var camera model.Camera

	if err := ctx.ShouldBindWith(&camera, binding.JSON); err == nil {
		c.CameraService.CreateOrUpdateCamera(configurationID, &camera)
		ctx.JSON(http.StatusAccepted, camera)
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

func (c Camera) deleteCamera(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	cameraID := ctx.Param("cameraId")

	if configuration := c.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	err := c.CameraService.DeleteCamera(cameraID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Deleted")
}

func (c Camera) toggleCameraAlarm(ctx *gin.Context) {

	configurationID := ctx.Param("id")
	cameraID := ctx.Param("cameraId")
	cameraIDInt, err := strconv.Atoi(cameraID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}
	statusString := ctx.DefaultQuery("status", "1")
	status, err := strconv.Atoi(statusString)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": err.Error()})
		return
	}

	if configuration := c.checkConfiguration(configurationID, ctx); configuration == nil {
		return
	}

	err = c.CameraService.ToggleCameraAlarm(uint(cameraIDInt), status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, "Done")
}

func (c Camera) checkConfiguration(configurationID string, ctx *gin.Context) *model.Configuration {
	configuration, err := c.ConfigurationService.GetConfiguration(configurationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": err.Error()})
		return nil
	}
	return configuration
}
