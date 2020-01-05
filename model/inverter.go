package model

import "github.com/jinzhu/gorm"

type Inverter struct {
	gorm.Model
	Name                       string `binding:"required"`
	Host                       string `binding:"required"`
	Port                       uint   `binding:"required"`
	ConfigurationID            uint
}
