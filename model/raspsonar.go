package model

import "github.com/jinzhu/gorm"

type Raspsonar struct {
	gorm.Model
	Host                          string  `binding:"required"`
	Port                          uint    `binding:"required"`
	SonarName                     string  `binding:"required"`
	RelayName                     string  `binding:"required"`
	DistanceThreshold             float64 `binding:"required"`
	AutoPowerOffDistanceThreshold float64 `binding:"required"`
	ConfigurationID               uint
}
