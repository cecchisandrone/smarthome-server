package service

import (
	"errors"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	"github.com/cecchisandrone/smarthome-server/slack"
	"github.com/spf13/viper"
)

type Alarm struct {
	SchedulerManager     *scheduler.SchedulerManager `inject:""`
	ConfigurationService *Configuration              `inject:""`
	NotificationService  *Notification               `inject:""`
	CameraService        *Camera                     `inject:""`
	SlackClient          slack.Client
	AlarmStatus          bool
	Cameras              []string
}

func (a *Alarm) Init() {
	a.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("alarm.locationChangeCheckIntervalSeconds")), a.checkLocation)
	configuration := a.ConfigurationService.GetCurrent()
	a.SlackClient = slack.Client{configuration.Slack}
	a.AlarmStatus = false
}

func (a *Alarm) ToggleAlarm(configuration model.Configuration, status int) ([]string, error) {

	errorString := ""
	a.Cameras = []string{}

	for _, c := range configuration.Cameras {
		if c.AlarmEnabled == true {
			err := a.CameraService.ToggleCameraAlarm(c.ID, status)
			if err != nil {
				errorString = errorString + "\n" + err.Error()
			} else {
				a.Cameras = append(a.Cameras, c.Name)
			}
		}
	}

	a.AlarmStatus = !(status == 0)

	if errorString == "" {
		return a.Cameras, nil
	}

	return a.Cameras, errors.New(errorString)
}

func (a *Alarm) checkLocation() {

}
