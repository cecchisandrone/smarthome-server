package model

import "github.com/jinzhu/gorm"

type CameraConfiguration struct {
	gorm.Model
	Name            string
	Host            string
	Username        string
	Password        string
	Enabled         bool
	AlarmEnabled    bool
	ConfigurationID uint
}
