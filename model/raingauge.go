package model

import "github.com/jinzhu/gorm"

type RainGauge struct {
	gorm.Model
	Host            string `binding:"required"`
	Port            uint   `binding:"required"`
	ConfigurationID uint
}
