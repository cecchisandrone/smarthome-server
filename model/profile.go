package model

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model
	Name            string
	Surname         string
	Password        string
	ConfigurationID uint
}
