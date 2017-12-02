package model

import "github.com/jinzhu/gorm"

type Configuration struct {
	gorm.Model
	Name        string
	Profile     Profile
	Gate        Gate
	Raspsonar   Raspsonar
	Temperature Temperature
	Slack       Slack
	Cameras     []Camera
}
