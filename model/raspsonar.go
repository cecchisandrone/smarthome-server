package model

import "github.com/jinzhu/gorm"

type Raspsonar struct {
	gorm.Model
	Host                          string  `binding:"required"`
	Port                          uint    `binding:"required"`
	SonarIndex                    uint    `binding:"required"`
	RelayIndex                    uint    `binding:"required"`
	DistanceThreshold             float32 `binding:"required"`
	AutoPowerOffDistanceThreshold float32 `binding:"required"`
	ConfigurationID               uint
}