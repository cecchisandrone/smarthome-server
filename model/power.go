package model

import "github.com/jinzhu/gorm"

type PowerMeter struct {
	gorm.Model
	Host             string  `binding:"required"`
	Port             uint    `binding:"required"`
	Voltage          float64 `binding:"required"`
	AdjustmentFactor float64 `binding:"required"`
	ConfigurationID  uint
}
