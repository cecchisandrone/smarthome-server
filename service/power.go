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

type PowerMeterMetrics struct {
	Power   float64 `json:"power"`
	Current float64 `json:"current"`
}

type PowerMeter struct {
	ScheduledMeasurements map[time.Time]float64
	SchedulerManager      *scheduler.SchedulerManager `inject:""`
	ConfigurationService  *Configuration              `inject:""`
	MaxMeasurements       int
	InfluxdbClient        *influxdb.Client `inject:""`
}

func (p *PowerMeter) Init() {
	p.ScheduledMeasurements = make(map[time.Time]float64)
	p.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("power.intervalSeconds")), p.ScheduledMeasurement)
	p.MaxMeasurements = viper.GetInt("power.maxMeasurements")
}

func (p *PowerMeter) GetLast(configuration model.Configuration) (time.Time, float64, error) {
	result := &PowerMeterMetrics{}
	_, err := resty.R().SetResult(&result).Get(getPowerMeterUrl(configuration))
	if err == nil {
		value := result.Current * configuration.PowerMeter.AdjustmentFactor * configuration.PowerMeter.Voltage
		return time.Now(), value, err
	} else {
		log.Error("Unable to fetch power measurement. Reason:", err)
		return time.Now(), 0, err
	}
}

func (p *PowerMeter) GetScheduledMeasurements() *map[time.Time]float64 {
	return &p.ScheduledMeasurements
}

func (p *PowerMeter) ScheduledMeasurement() {

	result := &PowerMeterMetrics{}
	configuration := p.ConfigurationService.GetCurrent()
	_, err := resty.R().SetResult(&result).Get(getPowerMeterUrl(configuration))
	if err == nil {
		value := result.Current * configuration.PowerMeter.AdjustmentFactor * configuration.PowerMeter.Voltage
		p.ScheduledMeasurements[time.Now()] = value
		log.Info("Scheduled power measurement: ", value)

		// Remove old measurements
		if len(p.ScheduledMeasurements) > p.MaxMeasurements {
			keys := make([]time.Time, 0)
			for key := range p.ScheduledMeasurements {
				keys = append(keys, key)
			}
			sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
			for _, key := range keys {
				if len(p.ScheduledMeasurements) > p.MaxMeasurements {
					delete(p.ScheduledMeasurements, key)
				}
			}
		}

		// Send data to influxdb
		point := client.Point{
			Measurement: "power",
			Tags:        nil,
			Fields: map[string]interface{}{
				"value": value,
			},
			Time: time.Now(),
		}
		p.InfluxdbClient.AddPoint(point)

	} else {
		log.Error("Unable to fetch power measurement. Reason:", err)
	}
}

func getPowerMeterUrl(configuration model.Configuration) string {
	return "http://" + configuration.PowerMeter.Host + ":" + strconv.FormatUint(uint64(configuration.PowerMeter.Port), 10) + "/metrics"
}
