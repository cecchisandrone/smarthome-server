package model

import "github.com/jinzhu/gorm"

type WellPump struct {
	gorm.Model
	Name                       string `binding:"required"`
	Host                       string `binding:"required"`
	Port                       uint   `binding:"required"`
	ActivationIntervals        string
	AutomaticActivationEnabled bool
	ManuallyActivated          bool
	ConfigurationID            uint
}
