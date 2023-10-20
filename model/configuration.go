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
	Alarm       Alarm
	Cameras     []Camera
	WellPumps   []WellPump
	RainGauge   RainGauge
	Humidity    Humidity
	Heater      Heater
	Inverters   []Inverter
	PowerMeter  PowerMeter
	Relays      []Relay
}
