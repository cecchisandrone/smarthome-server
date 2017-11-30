package model

import "github.com/jinzhu/gorm"

type Gate struct {
	gorm.Model
	Host            string `binding:"required"`
	Port            uint   `binding:"required"`
	Duration        float32 `binding:"required"`
	ConfigurationID uint
}
