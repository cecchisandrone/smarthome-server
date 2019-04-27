package service

import (
	"errors"
	"fmt"
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

type WellPump struct {
	Db                   *gorm.DB                    `inject:""`
	SchedulerManager     *scheduler.SchedulerManager `inject:""`
	ConfigurationService *Configuration              `inject:""`
	NotificationService  *Notification               `inject:""`
	RainGauge            *RainGauge                  `inject:""`
}

func (w *WellPump) Init() {
	w.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("well-pump.intervalSeconds")), w.ScheduledActivation)
}

func (w *WellPump) GetWellPumps(configurationID string) []model.WellPump {

	var WellPumps []model.WellPump
	w.Db.Where("configuration_id = ?", configurationID).Find(&WellPumps)

	return WellPumps
}

func (w *WellPump) CreateOrUpdateWellPump(configurationID string, wellPump *model.WellPump) {

	id, err := strconv.ParseUint(configurationID, 10, 32)
	if err == nil {
		wellPump.ConfigurationID = uint(id)
		w.Db.Save(&wellPump)
	}
}

func (w *WellPump) GetWellPump(wellPumpID string) (*model.WellPump, error) {

	var wellPump model.WellPump
	w.Db.First(&wellPump, wellPumpID)
	if wellPump.ID == 0 {
		return nil, errors.New("Can't find WellPump with ID " + string(wellPumpID))
	}
	return &wellPump, nil
}

func (w *WellPump) DeleteWellPump(wellPumpID string) error {

	var wellPump model.WellPump
	w.Db.First(&wellPump, wellPumpID)
	if wellPump.ID == 0 {
		return errors.New("Can't find WellPump with ID " + string(wellPumpID))
	}
	w.Db.Unscoped().Delete(&wellPump)
	return nil
}

func (w *WellPump) GetRelay(wellPump *model.WellPump) (int, error) {
	result := make(map[string]int)
	resp, err := resty.R().SetResult(&result).Get(getWellPumpUrl(wellPump))
	if err == nil && resp.StatusCode() == 200 {
		return result["status"], nil
	} else {
		return -1, errors.New("unable to get relay status")
	}
}

func (w *WellPump) ToggleRelay(wellPump *model.WellPump, status int, manuallyActivated bool) error {
	body := map[string]interface{}{"status": status}
	resp, err := resty.R().SetBody(body).Put(getWellPumpUrl(wellPump))
	if err != nil || resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("unable to toggle well pump %s to %d", wellPump.Name, status))
	}
	wellPump.ManuallyActivated = manuallyActivated
	w.Db.Save(&wellPump)

	return nil
}

func getWellPumpUrl(wellPump *model.WellPump) string {

	host := wellPump.Host
	port := wellPump.Port
	return fmt.Sprintf("http://%s:%d/relay", host, port)
}

func (w *WellPump) ScheduledActivation() {

	currentTime := time.Now()
	configuration := w.ConfigurationService.GetCurrent()

	for _, wellPump := range configuration.WellPumps {

		if !wellPump.ManuallyActivated && wellPump.AutomaticActivationEnabled {
			log.Info("Checking scheduled well pump " + wellPump.Name + " activation")
			status, err := w.GetRelay(&wellPump)
			if err != nil {
				w.NotificationService.SendSlackMessage(slack.AlarmChannel, "Unable to check status for well pump "+wellPump.Name)
				continue
			}

			// Skip pump schedule considering rainfall
			_, rainfall := w.RainGauge.GetLast24hTotal()
			if rainfall >= wellPump.RainfallThreshold {
				log.Info("Skipping well pump " + wellPump.Name + " schedule considering rainfall. Threshold: " + strconv.FormatFloat(wellPump.RainfallThreshold, 'f', 2, 64) + " - Rainfall: " + strconv.FormatFloat(rainfall, 'f', 2, 64))
				continue
			}

			timeIntervals, err := utils.ParseTimeIntervals(wellPump.ActivationIntervals)
			if err == nil {
				powerOffMatches := 0
				for label, startEndTimes := range timeIntervals {
					startTime := startEndTimes[0]
					endTime := startEndTimes[1]

					// Well pump currently off, check if it should be powered on
					if status == 0 && startTime.Before(currentTime) && endTime.After(currentTime) {

						log.Info("Turning on well pump " + wellPump.Name + " for interval " + label)
						err := w.ToggleRelay(&wellPump, 1, false)
						if err != nil {
							w.NotificationService.SendSlackMessage(slack.AlarmChannel, "Unable to turn on well pump "+wellPump.Name+" for interval "+label)
						}
						break
					}

					// Well pump currently on, check if it should be powered off
					if status == 1 && (startTime.After(currentTime) || endTime.Before(currentTime)) {
						powerOffMatches++
					}
				}
				// Well pump currently on, check if it should be powered off
				if len(timeIntervals) == powerOffMatches {
					log.Info("Turning off well pump " + wellPump.Name + ". We are out of interval/s " + wellPump.ActivationIntervals)
					err := w.ToggleRelay(&wellPump, 0, false)
					if err != nil {
						w.NotificationService.SendSlackMessage(slack.AlarmChannel, "Unable to turn off well pump "+wellPump.Name)
					}
				}
			}
		}
	}

}
