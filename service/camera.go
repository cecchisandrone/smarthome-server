package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/jinzhu/gorm"
	"gopkg.in/resty.v1"
)

type Camera struct {
	Db *gorm.DB `inject:""`
}

func (c Camera) Init() {

}

func (c Camera) GetCameras(configurationID string) []model.Camera {

	var Cameras []model.Camera
	c.Db.Where("configuration_id = ?", configurationID).Find(&Cameras)

	for i := range Cameras {
		generateUrl(&Cameras[i])
	}

	return Cameras
}

func (c Camera) CreateOrUpdateCamera(configurationID string, camera *model.Camera) {

	id, err := strconv.ParseUint(configurationID, 10, 32)
	if err == nil {
		camera.ConfigurationID = uint(id)
		c.Db.Save(&camera)
	}
	generateUrl(camera)
}

func (c Camera) GetCamera(cameraID string) (*model.Camera, error) {

	var camera model.Camera
	c.Db.First(&camera, cameraID)
	if camera.ID == 0 {
		return nil, errors.New("Can't find Camera with ID " + string(cameraID))
	}
	generateUrl(&camera)
	return &camera, nil
}

func (c Camera) DeleteCamera(cameraID string) error {

	var camera model.Camera
	c.Db.First(&camera, cameraID)
	if camera.ID == 0 {
		return errors.New("Can't find Camera with ID " + string(cameraID))
	}
	c.Db.Unscoped().Delete(&camera)
	return nil
}

func (c Camera) ToggleCameraAlarm(cameraID uint, status int) error {
	var camera model.Camera
	c.Db.First(&camera, cameraID)
	if camera.ID == 0 {
		return errors.New("Can't find Camera with ID " + string(cameraID))
	}

	if camera.Type != model.Foscam {
		return errors.New("Can't toggle alarm on camera type " + string(camera.Type))
	}

	url := fmt.Sprintf("http://%s:%d/set_alarm.cgi?mail=%d&user=%s&pwd=%s", camera.Host, camera.Port, status, camera.Username, camera.Password)
	_, err := resty.R().Get(url)
	return err
}

func generateUrl(camera *model.Camera) {
	switch camera.Type {
	case model.Foscam:
		camera.Url = fmt.Sprintf("http://%s:%d/videostream.cgi?user=%s&pwd=%s", camera.Host, camera.Port, camera.Username, camera.Password)
	case model.ADJ:
		camera.Url = fmt.Sprintf("http://%s:%d/videostream.cgi?user=%s&pwd=%s", camera.Host, camera.Port, camera.Username, camera.Password)
	case model.Microcam:
		camera.Url = fmt.Sprintf("http://%s:%d/media/?action=stream&user=%s&pwd=%s", camera.Host, camera.Port, camera.Username, camera.Password)
	case model.SV3C:
		camera.Url = fmt.Sprintf("http://%s:%d/sv3c/tmpfs/snap.jpg", camera.Host, camera.Port)
	}
}
