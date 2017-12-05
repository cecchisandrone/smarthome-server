package service

import (
	"fmt"
	"github.com/cecchisandrone/smarthome-server/model"
	"gopkg.in/resty.v1"
)

type Gate struct {
	ConfigurationService *Configuration `inject:""`
}

func (g *Gate) Init() {
}

func (g *Gate) Open(configuration model.Configuration) error {
	_, err := resty.R().Post(getGateUrl(configuration))
	return err
}

func getGateUrl(configuration model.Configuration) string {

	host := configuration.Gate.Host
	port := configuration.Gate.Port
	duration := configuration.Gate.Duration
	return fmt.Sprintf("http://%s:%d/toggle-relay?duration=%d", host, port, duration)
}
