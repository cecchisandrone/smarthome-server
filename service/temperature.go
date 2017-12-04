package service

import (
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"strconv"
	"time"
)

type Temperature struct {
	ScheduledMeasurements map[time.Time]float64
	SchedulerManager      *scheduler.SchedulerManager `inject:""`
	ConfigurationService  *Configuration              `inject:""`
	MaxMeasurements       int
}

func (t *Temperature) Init() {
	t.ScheduledMeasurements = make(map[time.Time]float64)
	t.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("temperature.intervalSeconds")), t.ScheduledMeasurement)
	t.MaxMeasurements = viper.GetInt("temperature.maxMeasurements")
}

func (t *Temperature) GetLast(configuration model.Configuration) (time.Time, float64) {
	resp, err := resty.R().Get(getTemperatureUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		return time.Now(), value
	} else {
		log.Error("Unable to fetch temperature measurement. Reason:", err)
		return time.Now(), 0
	}
}

func (t *Temperature) GetScheduledMeasurements() *map[time.Time]float64 {
	return &t.ScheduledMeasurements
}

func (t *Temperature) ScheduledMeasurement() {

	configuration := t.ConfigurationService.GetCurrent()
	resp, err := resty.R().Get(getTemperatureUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		t.ScheduledMeasurements[time.Now()] = value
		log.Info("Scheduled temperature measurement:", value)

		// Remove old measurements
		index := 0
		if len(t.ScheduledMeasurements) > t.MaxMeasurements {
			for key := range t.ScheduledMeasurements {
				index++
				if index >= t.MaxMeasurements {
					delete(t.ScheduledMeasurements, key)
				}
			}
		}
	} else {
		log.Error("Unable to fetch temperature measurement. Reason:", err)
	}
}

func getTemperatureUrl(configuration model.Configuration) string {
	return "http://" + configuration.Temperature.Host + ":" + strconv.FormatUint(uint64(configuration.Temperature.Port), 10) + "/temp"
}
