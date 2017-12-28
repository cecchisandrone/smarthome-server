package model

import (
	"github.com/jinzhu/gorm"
)

type Alarm struct {
	gorm.Model
	AutomaticAlarmActivation bool
	ConfigurationID          uint
}
