package service

import (
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"sort"
	"strconv"
	"time"
)

type RainGauge struct {
	ScheduledMeasurements map[time.Time]float64
	SchedulerManager      *scheduler.SchedulerManager `inject:""`
	ConfigurationService  *Configuration              `inject:""`
	MaxMeasurements       int
}

func (r *RainGauge) Init() {
	r.ScheduledMeasurements = make(map[time.Time]float64)
	r.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("rain-gauge.intervalSeconds")), r.ScheduledMeasurement)
	r.MaxMeasurements = viper.GetInt("rain-gauge.maxMeasurements")
}

func (r *RainGauge) GetLast24hTotal() (time.Time, float64) {
	startTime := time.Now().AddDate(0, 0, -1)
	sum := 0.
	for key, value := range r.ScheduledMeasurements {
		if key.After(startTime) {
			sum += value
		}
	}
	return time.Now(), sum
}

func (r *RainGauge) GetScheduledMeasurements() *map[time.Time]float64 {
	return &r.ScheduledMeasurements
}

func (r *RainGauge) ScheduledMeasurement() {

	configuration := r.ConfigurationService.GetCurrent()
	resp, err := resty.R().Get(getRainGaugeUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		r.ScheduledMeasurements[time.Now()] = value
		log.Info("Scheduled rain gauge measurement: ", value)

		// Remove old measurements
		if len(r.ScheduledMeasurements) > r.MaxMeasurements {
			keys := make([]time.Time, 0)
			for key := range r.ScheduledMeasurements {
				keys = append(keys, key)
			}
			sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
			for _, key := range keys {
				if len(r.ScheduledMeasurements) > r.MaxMeasurements {
					delete(r.ScheduledMeasurements, key)
				}
			}
		}
	} else {
		log.Error("Unable to fetch rain gauge measurement. Reason:", err)
	}
}

func getRainGaugeUrl(configuration model.Configuration) string {
	return "http://" + configuration.RainGauge.Host + ":" + strconv.FormatUint(uint64(configuration.RainGauge.Port), 10) + "/"
}
