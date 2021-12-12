package model

import "github.com/jinzhu/gorm"

type Heater struct {
	gorm.Model
	Host            string `binding:"required"`
	Port            uint   `binding:"required"`
	ConfigurationID uint
}
