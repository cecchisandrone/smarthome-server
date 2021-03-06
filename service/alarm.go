package service

import (
	"errors"
	"strings"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	"github.com/cecchisandrone/smarthome-server/slack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type LocationStatus string

const (
	Entered LocationStatus = "entered"
	Exited  LocationStatus = "exited"
)

type Alarm struct {
	SchedulerManager     *scheduler.SchedulerManager `inject:""`
	ConfigurationService *Configuration              `inject:""`
	NotificationService  *Notification               `inject:""`
	CameraService        *Camera                     `inject:""`
	SlackClient          slack.Client
	AlarmStatus          bool
	Cameras              []string
	LocationStatus       map[string]LocationStatus
}

func (a *Alarm) Init() {
	a.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("alarm.locationChangeCheckIntervalSeconds")), a.checkLocationStatus)
	a.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("alarm.automaticAlarmToggleIntervalSeconds")), a.automaticAlarmToggle)
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

func (a *Alarm) automaticAlarmToggle() {

	configuration := a.ConfigurationService.GetCurrent()
	if configuration.Alarm.AutomaticAlarmActivation {

		log.Info("Automatic alarm activation is enabled. Users status: ", a.LocationStatus)

		if len(a.LocationStatus) > 0 {
			for _, status := range a.LocationStatus {
				if status == Entered {
					// At least one entered, disable alarm
					a.ToggleAlarm(configuration, 0)
					return
				}
			}

			a.ToggleAlarm(configuration, 1)
		}
	}
}

func (a *Alarm) checkLocationStatus() {

	a.LocationStatus = make(map[string]LocationStatus)

	configuration := a.ConfigurationService.GetCurrent()

	users := configuration.Slack.GetLocationChangeUsersArray()

	history, err := a.SlackClient.GetLocationChangeChannelHistory(configuration.Slack.LocationChangeChannel)
	if err == nil && history != nil {
		for _, message := range history.Messages {
			if len(message.Attachments) != 0 {
				text := message.Attachments[0].Text
				tokens := strings.Split(text, " ")
				for _, user := range users {
					_, present := a.LocationStatus[user]
					if tokens[0] == user && !present {
						a.LocationStatus[user] = LocationStatus(tokens[1])
						break
					}
				}
			}
			if len(a.LocationStatus) == len(users) {
				break
			}
		}
	} else {
		log.Error("Error while checking location status. Reason: ", err.Error())
	}
}
