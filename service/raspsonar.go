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
	"math"
)

type Raspsonar struct {
	ScheduledMeasurements    map[time.Time]float64
	SchedulerManager         *scheduler.SchedulerManager `inject:""`
	ConfigurationService     *Configuration              `inject:""`
	NotificationService      *Notification               `inject:""`
	MaxMeasurements          int
	RelayStatus              bool
	RelayActivationTimestamp time.Time
	Configuration            model.Configuration
	lastMeasure              float64
}

func (r *Raspsonar) Init() {
	r.ScheduledMeasurements = make(map[time.Time]float64)
	r.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("raspsonar.intervalSeconds")), r.ScheduledMeasurement)
	r.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("raspsonar.autoToggleRelayIntervalSeconds")), r.autoToggleRelay)
	r.MaxMeasurements = viper.GetInt("raspsonar.maxMeasurements")
	r.RelayStatus = false
	r.Configuration = r.ConfigurationService.GetCurrent()
}

func (r *Raspsonar) GetLast(configuration model.Configuration) (time.Time, float64, error) {
	resp, err := resty.R().Get(getDistanceUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)

		// Average value with previous (reduce noise)
		if r.lastMeasure != 0 {
			// Skip values too distant from previous, should be an error
			if math.Abs(value-r.lastMeasure) > 5 {
				log.Warn("Ignoring raspsonar value " + strconv.FormatFloat(value, 'f', 2, 64))
				value = r.lastMeasure
			}
			r.lastMeasure = value*0.3 + r.lastMeasure*0.7
		} else {
			r.lastMeasure = value
		}

		return time.Now(), r.lastMeasure, nil
	} else {
		log.Error("Unable to fetch raspsonar measurement. Reason:", err)
		return time.Now(), 0, errors.New("Unable to fetch raspsonar measurement")
	}
}

func (r *Raspsonar) ToggleRelay(configuration model.Configuration, status bool) error {
	statusInt := 1
	if status {
		statusInt = 0
	}
	_, err := resty.R().Put(getToggleRelayUrl(configuration, statusInt))
	if err != nil {
		log.Error("Unable to toggle relay. Reason:", err)
		return err
	}
	r.RelayStatus = status
	if r.RelayStatus == true {
		r.RelayActivationTimestamp = time.Now()
	}
	return nil
}

func (r *Raspsonar) GetScheduledMeasurements() *map[time.Time]float64 {
	return &r.ScheduledMeasurements
}

func (r *Raspsonar) ScheduledMeasurement() {

	timestamp, value, err := r.GetLast(r.Configuration)
	if err == nil {
		r.ScheduledMeasurements[timestamp] = value
		log.Info("Scheduled raspsonar measurement: " + strconv.FormatFloat(value, 'f', 2, 64))

		if value < r.Configuration.Raspsonar.DistanceThreshold {
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

func (r *Raspsonar) autoToggleRelay() {
	if r.RelayStatus {
		_, distance, err := r.GetLast(r.Configuration)
		if err == nil {

			log.Info("Checking if relay should be put off. Distance: " + strconv.FormatFloat(distance, 'f', 2, 64))

			// Toggle relay off is threshold is trespassed
			if distance > r.Configuration.Raspsonar.AutoPowerOffDistanceThreshold {
				log.Info("Toggling relay off...threshold is trespassed")
				r.ToggleRelay(r.Configuration, false)
				r.NotificationService.SendSlackMessage(slack.AlarmChannel, "Auto power off distance threshold ("+strconv.FormatFloat(r.Configuration.Raspsonar.AutoPowerOffDistanceThreshold, 'f', 2, 64)+") trespassed. Powering off the pump")

			}
		} else {
			log.Error("Error while getting auto power off measurement. Reason: ", err.Error())
		}
	}
}
