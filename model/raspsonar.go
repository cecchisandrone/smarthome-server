package model

import "github.com/jinzhu/gorm"

type Raspsonar struct {
	gorm.Model
	Url                           string  `binding:"required"`
	SonarIndex                    uint    `binding:"required"`
	RelayIndex                    uint    `binding:"required"`
	DistanceThreshold             float32 `binding:"required"`
	AutoPowerOffDistanceThreshold float32 `binding:"required"`
	ConfigurationID               uint
}
