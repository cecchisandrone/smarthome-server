package model

import (
	"database/sql/driver"
	"github.com/jinzhu/gorm"
)

type CameraType string

func (u *CameraType) Scan(value interface{}) error {
	*u = CameraType(value.([]byte))
	return nil
}

func (u CameraType) Value() (driver.Value, error) {
	return string(u), nil
}

const (
	Foscam   CameraType = "foscam"
	ADJ      CameraType = "adj"
	Microcam CameraType = "microcam"
)

type Camera struct {
	gorm.Model
	Name            string
	Type            CameraType `sql:"not null;type:ENUM('foscam', 'adj', 'microcam')"`
	Host            string
	Port            uint
	Url             string `gorm:"-"`
	Username        string
	Password        string
	Enabled         bool
	AlarmEnabled    bool
	ConfigurationID uint
}