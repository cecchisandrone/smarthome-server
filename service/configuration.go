package service

import (
	"errors"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/jinzhu/gorm"
)

type Configuration struct {
	Db *gorm.DB `inject:""`
}

func (c *Configuration) GetConfigurations() []model.Configuration {

	var configurations []model.Configuration
	c.Db.Find(&configurations)
	return configurations
}

func (c *Configuration) CreateConfiguration(configuration *model.Configuration) {
	c.Db.Save(configuration)
}

func (c *Configuration) GetConfiguration(configurationID string) (*model.Configuration, error) {

	var configuration model.Configuration
	c.Db.Preload("Profile").Preload("CameraConfigurations").First(&configuration, configurationID)
	if configuration.ID == 0 {
		return nil, errors.New("Can't find configuration with ID " + string(configurationID))
	}
	return &configuration, nil
}

func (c *Configuration) DeleteConfiguration(configurationID string) error {

	var configuration model.Configuration
	c.Db.First(&configuration, configurationID)
	if configuration.ID == 0 {
		return errors.New("Can't find configuration with ID " + string(configurationID))
	}
	c.Db.Unscoped().Delete(&configuration)
	return nil
}
