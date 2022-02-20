package service

import (
	"errors"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/jinzhu/gorm"
)

type Configuration struct {
	Db *gorm.DB `inject:""`
}

func (c *Configuration) Init() {

}

func (c *Configuration) GetConfigurations() []model.Configuration {

	var configurations []model.Configuration
	c.Db.Preload("Profile").Preload("Cameras").Preload("Gate").Preload("Temperature").Preload("Slack").Preload("Raspsonar").Preload("Alarm").Preload("WellPumps").Preload("RainGauge").Preload("Humidity").Preload("Inverters").Preload("Heater").Preload("PowerMeter").Find(&configurations)
	return configurations
}

func (c *Configuration) GetCurrent() model.Configuration {
	var configuration model.Configuration
	c.Db.Preload("Profile").Preload("Cameras").Preload("Gate").Preload("Temperature").Preload("Slack").Preload("Raspsonar").Preload("Alarm").Preload("WellPumps").Preload("Cameras").Preload("RainGauge").Preload("Humidity").Preload("Inverters").Preload("Heater").Preload("PowerMeter").Last(&configuration)
	return configuration
}

func (c *Configuration) GetConfiguration(configurationID string) (*model.Configuration, error) {

	var configuration model.Configuration
	c.Db.Preload("Profile").Preload("Cameras").Preload("Gate").Preload("Temperature").Preload("Slack").Preload("Raspsonar").Preload("Alarm").Preload("WellPumps").Preload("Cameras").Preload("RainGauge").Preload("Humidity").Preload("Inverters").Preload("Heater").Preload("PowerMeter").First(&configuration, configurationID)
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
func (c *Configuration) CreateOrUpdateConfiguration(configuration *model.Configuration) {
	c.Db.Save(&configuration)
}
