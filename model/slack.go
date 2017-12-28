package model

import (
	"github.com/jinzhu/gorm"
	"strings"
)

type Slack struct {
	gorm.Model
	NotificationChannel   string `binding:"required"`
	LocationChangeChannel string `binding:"required"`
	Token                 string `binding:"required"`
	LocationChangeUsers   string `binding:"required"`
	ConfigurationID       uint
}

func (s Slack) GetLocationChangeUsersArray() []string {
	return strings.Split(s.LocationChangeUsers, ";")
}
