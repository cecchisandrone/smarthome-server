package model

import "github.com/jinzhu/gorm"

type Rental struct {
	gorm.Model
	Url             string `binding:"required"`
	ConfigurationID uint
}
