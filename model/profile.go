package model

import "github.com/jinzhu/gorm"

type Profile struct {
	gorm.Model
	FirstName       string `binding:"required"`
	LastName        string `binding:"required"`
	Password        string
	Username        string `binding:"required"`
	ConfigurationID uint
}
