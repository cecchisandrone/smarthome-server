package service

import (
	"github.com/cecchisandrone/smarthome-server/influxdb"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	client "github.com/influxdata/influxdb1-client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"sort"
	"strconv"
	"time"
)

type Heating struct {
	ScheduledMeasurements map[time.Time]float64
	SchedulerManager      *scheduler.SchedulerManager `inject:""`
	ConfigurationService  *Configuration              `inject:""`
	MaxMeasurements       int
	InfluxdbClient        *influxdb.Client `inject:""`
}

func (h *Heating) Init() {
	h.ScheduledMeasurements = make(map[time.Time]float64)
	h.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("heating.intervalSeconds")), h.ScheduledMeasurement)
	h.MaxMeasurements = viper.GetInt("heating.maxMeasurements")
}

func (h *Heating) GetLast(configuration model.Configuration) (time.Time, float64, error) {
	resp, err := resty.R().Get(getHeatingUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		return time.Now(), value, err
	} else {
		log.Error("Unable to fetch heating measurement. Reason:", err)
		return time.Now(), 0, err
	}
}

func (h *Heating) GetScheduledMeasurements() *map[time.Time]float64 {
	return &h.ScheduledMeasurements
}

func (h *Heating) ScheduledMeasurement() {

	configuration := h.ConfigurationService.GetCurrent()
	resp, err := resty.R().Get(getTemperatureUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		h.ScheduledMeasurements[time.Now()] = value
		log.Info("Scheduled heating measurement: ", value)

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

		// Send data to influxdb
		point := client.Point{
			Measurement: "temperature",
			Tags: map[string]string{
				"location": "heating",
			},
			Fields: map[string]interface{}{
				"value": value,
			},
			Time: time.Now(),
		}
		h.InfluxdbClient.AddPoint(point)

	} else {
		log.Error("Unable to fetch heating measurement. Reason:", err)
	}
}

func getHeatingUrl(configuration model.Configuration) string {
	return "http://" + configuration.Temperature.Host + ":" + strconv.FormatUint(uint64(configuration.Temperature.Port), 10) + "/temp"
}
