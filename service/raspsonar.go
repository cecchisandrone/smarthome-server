package service

import (
	"fmt"
	"strconv"
	"time"

	"errors"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	"github.com/cecchisandrone/smarthome-server/slack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
)

type Raspsonar struct {
	ScheduledMeasurements    map[time.Time]float64
	SchedulerManager         *scheduler.SchedulerManager `inject:""`
	ConfigurationService     *Configuration              `inject:""`
	NotificationService      *Notification               `inject:""`
	MaxMeasurements          int
	RelayStatus              bool
	RelayActivationTimestamp time.Time
}

func (r *Raspsonar) Init() {
	r.ScheduledMeasurements = make(map[time.Time]float64)
	r.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("raspsonar.intervalSeconds")), r.ScheduledMeasurement)
	r.MaxMeasurements = viper.GetInt("raspsonar.maxMeasurements")
	r.RelayStatus = false
}

func (r *Raspsonar) GetLast(configuration model.Configuration) (time.Time, float64, error) {
	resp, err := resty.R().Get(getDistanceUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		return time.Now(), value, nil
	} else {
		log.Error("Unable to fetch raspsonar measurement. Reason:", err)
		return time.Now(), 0, errors.New("Unable to fetch raspsonar measurement")
	}
}

func (r *Raspsonar) ToggleRelay(configuration model.Configuration, status int) error {
	_, err := resty.R().Put(getToggleRelayUrl(configuration, status))
	if err != nil {
		log.Error("Unable to toggle relay. Reason:", err)
		return err
	}
	r.RelayStatus = !(status == 0)
	if r.RelayStatus == true {
		r.RelayActivationTimestamp = time.Now()
	}
	return nil
}

func (r *Raspsonar) GetScheduledMeasurements() *map[time.Time]float64 {
	return &r.ScheduledMeasurements
}

func (r *Raspsonar) ScheduledMeasurement() {

	configuration := r.ConfigurationService.GetCurrent()
	resp, err := resty.R().Get(getDistanceUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		r.ScheduledMeasurements[time.Now()] = value
		log.Info("Scheduled raspsonar measurement:", value)

		if value < configuration.Raspsonar.DistanceThreshold {
			r.NotificationService.SendSlackMessage(slack.AlarmChannel, "Warning! Distance threshold has been trespassed. Value: "+strconv.FormatFloat(value, 'f', 2, 64))
		}

		// Remove old measurements
		index := 0
		if len(r.ScheduledMeasurements) > r.MaxMeasurements {
			for key := range r.ScheduledMeasurements {
				index++
				if index >= r.MaxMeasurements {
					delete(r.ScheduledMeasurements, key)
				}
			}
		}
	} else {
		log.Error("Unable to fetch raspsonar measurement. Reason:", err)
	}
}

func getDistanceUrl(configuration model.Configuration) string {

	host := configuration.Raspsonar.Host
	port := configuration.Raspsonar.Port
	name := configuration.Raspsonar.SonarName
	return fmt.Sprintf("http://%s:%d/devices/sonar/%s", host, port, name)
}

func getToggleRelayUrl(configuration model.Configuration, status int) string {

	host := configuration.Raspsonar.Host
	port := configuration.Raspsonar.Port
	name := configuration.Raspsonar.RelayName
	return fmt.Sprintf("http://%s:%d/devices/relay/%s?status=%d", host, port, name, status)
}
