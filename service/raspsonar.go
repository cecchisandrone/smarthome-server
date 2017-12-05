package service

import (
	"fmt"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"strconv"
	"time"
)

type Raspsonar struct {
	ScheduledMeasurements map[time.Time]float64
	SchedulerManager      *scheduler.SchedulerManager `inject:""`
	ConfigurationService  *Configuration              `inject:""`
	MaxMeasurements       int
}

func (r *Raspsonar) Init() {
	r.ScheduledMeasurements = make(map[time.Time]float64)
	r.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("raspsonar.intervalSeconds")), r.ScheduledMeasurement)
	r.MaxMeasurements = viper.GetInt("raspsonar.maxMeasurements")
}

func (r *Raspsonar) GetLast(configuration model.Configuration) (time.Time, float64) {
	resp, err := resty.R().Get(getRaspsonarUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		return time.Now(), value
	} else {
		log.Error("Unable to fetch raspsonar measuremenr. Reason:", err)
		return time.Now(), 0
	}
}

func (r *Raspsonar) GetScheduledMeasurements() *map[time.Time]float64 {
	return &r.ScheduledMeasurements
}

func (r *Raspsonar) ScheduledMeasurement() {

	configuration := r.ConfigurationService.GetCurrent()
	resp, err := resty.R().Get(getRaspsonarUrl(configuration))
	if err == nil {
		value, _ := strconv.ParseFloat(resp.String(), 64)
		r.ScheduledMeasurements[time.Now()] = value
		log.Info("Scheduled raspsonar measurement:", value)

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
		log.Error("Unable to fetch raspsonar measuremenr. Reason:", err)
	}
}

func getRaspsonarUrl(configuration model.Configuration) string {

	host := configuration.Raspsonar.Host
	port := configuration.Raspsonar.Port
	index := configuration.Raspsonar.SonarIndex
	measurements := viper.GetInt("raspsonar.measurements")
	return fmt.Sprintf("http://%s:%d/raspio/rest/sonar/%d/distance?measurements=%d", host, port, index, measurements)
}
