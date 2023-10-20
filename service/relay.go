package service

import (
	"errors"
	"fmt"
	"github.com/cecchisandrone/smarthome-server/dto"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	"github.com/cecchisandrone/smarthome-server/slack"
	"github.com/cecchisandrone/smarthome-server/utils"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"strconv"
	"time"
)

type Relay struct {
	Db                   *gorm.DB                    `inject:""`
	SchedulerManager     *scheduler.SchedulerManager `inject:""`
	ConfigurationService *Configuration              `inject:""`
	NotificationService  *Notification               `inject:""`
}

func (r *Relay) Init() {
	r.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("relay.intervalSeconds")), r.ScheduledActivation)
}

func (r *Relay) GetRelays(configurationID string) []model.Relay {

	var relays []model.Relay
	r.Db.Where("configuration_id = ?", configurationID).Find(&relays)

	return relays
}

func (r *Relay) CreateOrUpdateRelay(configurationID string, Relay *model.Relay) {

	id, err := strconv.ParseUint(configurationID, 10, 32)
	if err == nil {
		Relay.ConfigurationID = uint(id)
		r.Db.Save(&Relay)
	}
}

func (r *Relay) GetRelay(RelayID string) (*model.Relay, error) {

	var Relay model.Relay
	r.Db.First(&Relay, RelayID)
	if Relay.ID == 0 {
		return nil, errors.New("Can't find Relay with ID " + string(RelayID))
	}
	return &Relay, nil
}

func (r *Relay) DeleteRelay(RelayID string) error {

	var Relay model.Relay
	r.Db.First(&Relay, RelayID)
	if Relay.ID == 0 {
		return errors.New("Can't find Relay with ID " + string(RelayID))
	}
	r.Db.Unscoped().Delete(&Relay)
	return nil
}

func (r *Relay) GetRelayStatus(Relay *model.Relay) (map[int]bool, error) {

	type Pin struct {
		Pin    int  `json:"pin"`
		Status bool `json:"status"`
	}

	var result []Pin
	resp, err := resty.R().SetResult(&result).Get(getRelayUrl(Relay))
	if err == nil && resp.StatusCode() == 200 {
		pinMap := make(map[int]bool)
		for _, pinStatus := range result {
			pinMap[pinStatus.Pin] = pinStatus.Status
		}
		return pinMap, nil
	} else {
		return nil, errors.New("unable to get relay status")
	}
}

func (r *Relay) ToggleRelay(Relay *model.Relay, status []dto.Pin, manuallyActivated bool) error {

	resp, err := resty.R().SetBody(status).Put(getRelayUrl(Relay))
	if err != nil || resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("unable to toggle relay %s to %d", Relay.Name, status))
	}
	Relay.ManuallyActivated = manuallyActivated
	r.Db.Save(&Relay)

	return nil
}

func (r *Relay) ToggleAllPinsRelay(Relay *model.Relay, status bool) error {
	var body []dto.Pin

	for i := uint(0); i < Relay.Channels; i++ {
		pin := dto.Pin{
			Pin:    int(i),
			Status: status,
		}
		body = append(body, pin)
	}

	resp, err := resty.R().SetBody(body).Put(getRelayUrl(Relay))
	if err != nil || resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("unable to toggle relay %s to %d", Relay.Name, status))
	}
	r.Db.Save(&Relay)

	return nil
}

func getRelayUrl(Relay *model.Relay) string {

	host := Relay.Host
	port := Relay.Port
	return fmt.Sprintf("http://%s:%d/relay", host, port)
}

func (r *Relay) ScheduledActivation() {

	currentTime := time.Now()
	configuration := r.ConfigurationService.GetCurrent()

	for _, Relay := range configuration.Relays {

		if !Relay.ManuallyActivated && Relay.AutomaticActivationEnabled {
			log.Info("Checking scheduled relay " + Relay.Name + " activation")
			status, err := r.GetRelayStatus(&Relay)
			if err != nil {
				r.NotificationService.SendSlackMessage(slack.AlarmChannel, "Unable to check status for relay "+Relay.Name)
				continue
			}

			// range over status map to check if any pin is on
			for _, pinStatus := range status {

				timeIntervals, err := utils.ParseTimeIntervals(Relay.ActivationIntervals)
				if err == nil {
					powerOffMatches := 0
					for label, startEndTimes := range timeIntervals {
						startTime := startEndTimes[0]
						endTime := startEndTimes[1]

						// Relay currently off, check if it should be powered on
						if !pinStatus && startTime.Before(currentTime) && endTime.After(currentTime) {

							log.Info("Turning on relay " + Relay.Name + " for interval " + label)
							err := r.ToggleAllPinsRelay(&Relay, true)
							if err != nil {
								r.NotificationService.SendSlackMessage(slack.AlarmChannel, "Unable to turn on relay "+Relay.Name+" for interval "+label)
							}
							break
						}

						// Relay on, check if it should be powered off
						if !pinStatus && (startTime.After(currentTime) || endTime.Before(currentTime)) {
							powerOffMatches++
						}
					}
					// Relay currently on, check if it should be powered off
					if len(timeIntervals) == powerOffMatches {
						log.Info("Turning off relay " + Relay.Name + ". We are out of interval/s " + Relay.ActivationIntervals)
						err := r.ToggleAllPinsRelay(&Relay, false)
						if err != nil {
							r.NotificationService.SendSlackMessage(slack.AlarmChannel, "Unable to turn off relay "+Relay.Name)
						}
					}
				}
			}
		}
	}

}
