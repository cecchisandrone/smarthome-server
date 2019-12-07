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

type Temperature struct {
	ScheduledMeasurements map[time.Time]float64
	SchedulerManager      *scheduler.SchedulerManager `inject:""`
	ConfigurationService  *Configuration              `inject:""`
	MaxMeasurements       int
	InfluxdbClient  *influxdb.Client                  `inject:""`
}

func (t *Temperature) Init() {
	t.ScheduledMeasurements = make(map[time.Time]float64)
	t.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("temperature.intervalSeconds")), t.ScheduledMeasurement)
	t.MaxMeasurements = viper.GetInt("temperature.maxMeasurements")
}

func (t *Temperature) GetLast(configuration model.Configuration) (time.Time, float64, error) {
	resp, err := resty.R().Get(getTemperatureUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		return time.Now(), value, err
	} else {
		log.Error("Unable to fetch temperature measurement. Reason:", err)
		return time.Now(), 0, err
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
		log.Info("Scheduled temperature measurement: ", value)

		// Remove old measurements
		if len(t.ScheduledMeasurements) > t.MaxMeasurements {
			keys := make([]time.Time, 0)
			for key := range t.ScheduledMeasurements {
				keys = append(keys, key)
			}
			sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
			for _, key := range keys {
				if len(t.ScheduledMeasurements) > t.MaxMeasurements {
					delete(t.ScheduledMeasurements, key)
				}
			}
		}

		// Send data to influxdb
		point := client.Point{
			Measurement: "temperature",
			Tags: map[string]string{
				"location": "gate",
			},
			Fields: map[string]interface{}{
				"value": value,
			},
			Time:      time.Now(),
		}
		t.InfluxdbClient.AddPoint(point)

	} else {
		log.Error("Unable to fetch temperature measurement. Reason:", err)
	}
}

func getTemperatureUrl(configuration model.Configuration) string {
	return "http://" + configuration.Temperature.Host + ":" + strconv.FormatUint(uint64(configuration.Temperature.Port), 10) + "/temp"
}
