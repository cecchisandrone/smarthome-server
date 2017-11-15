package model

import "github.com/jinzhu/gorm"

type Configuration struct {
	gorm.Model
	Name                 string
	Profile              Profile
	CameraConfigurations []CameraConfiguration
}
