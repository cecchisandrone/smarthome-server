package service

import (
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"math"
	"sort"
	"strconv"
	"time"
)

type Humidity struct {
	ScheduledMeasurements     map[time.Time]float64
	SchedulerManager          *scheduler.SchedulerManager `inject:""`
	ConfigurationService      *Configuration              `inject:""`
	MaxMeasurements           int
	lastMeasure               float64
	wrongMeasurementThreshold float64
}

func (h *Humidity) Init() {
	h.ScheduledMeasurements = make(map[time.Time]float64)
	h.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("humidity.intervalSeconds")), h.ScheduledMeasurement)
	h.MaxMeasurements = viper.GetInt("humidity.maxMeasurements")
}

func (h *Humidity) GetLast(configuration model.Configuration) (time.Time, float64, error) {
	resp, err := resty.R().Get(getHumidityUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)

		// Average value with previous (reduce noise)
		if h.lastMeasure != 0 {
			// Skip values too distant from previous, should be an error
			if math.Abs(value-h.lastMeasure) > h.wrongMeasurementThreshold {
				log.Warn("Ignoring humidity value " + strconv.FormatFloat(value, 'f', 2, 64))
				value = h.lastMeasure
			}
			h.lastMeasure = value*0.3 + h.lastMeasure*0.7
		} else {
			h.lastMeasure = value
		}

		return time.Now(), h.lastMeasure, nil
	} else {
		log.Error("Unable to fetch humidity measurement. Reason:", err)
		return time.Now(), 0, err
	}
}

func (h *Humidity) GetScheduledMeasurements() *map[time.Time]float64 {
	return &h.ScheduledMeasurements
}

func (h *Humidity) ScheduledMeasurement() {

	configuration := h.ConfigurationService.GetCurrent()
	resp, err := resty.R().Get(getHumidityUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		h.ScheduledMeasurements[time.Now()] = value
		log.Info("Scheduled humidity measurement: ", value)

		// Remove old measurements
		if len(h.ScheduledMeasurements) > h.MaxMeasurements {
			keys := make([]time.Time, 0)
			for key := range h.ScheduledMeasurements {
				keys = append(keys, key)
			}
			sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
			for _, key := range keys {
				if len(h.ScheduledMeasurements) > h.MaxMeasurements {
					delete(h.ScheduledMeasurements, key)
				}
			}
		}
	} else {
		log.Error("Unable to fetch humidity measurement. Reason:", err)
	}
}

func getHumidityUrl(configuration model.Configuration) string {
	return "http://" + configuration.Humidity.Host + ":" + strconv.FormatUint(uint64(configuration.Humidity.Port), 10) + "/humidity"
}
