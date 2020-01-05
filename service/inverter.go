package service

import (
	"errors"
	"fmt"
	"github.com/cecchisandrone/smarthome-server/influxdb"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/cecchisandrone/smarthome-server/scheduler"
	client "github.com/influxdata/influxdb1-client"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/resty.v1"
	"strconv"
	"strings"
	"sync"
	"time"
)

type InverterMetrics struct {
	PowerPin1 float32 `json:"power_pin_1"`
	PowerPin2  float32 `json:"power_pin_2"`
	GridPowerReading float32 `json:"grid_power_reading"`
	Riso float32 `json:"riso"`
	InverterTemperature float32 `json:"inverter_temperature"`
	BoosterTemperature float32 `json:"booster_temperature"`
	DcAcConversionEfficiency float32 `json:"dc_ac_conversion_efficiency"`
	DailyEnergy float32 `json:"daily_energy"`
	WeeklyEnergy float32 `json:"weekly_energy"`
	MontlyEnergy float32 `json:"monthly_energy"`
	YearlyEnergy float32 `json:"yearly_energy"`
	PowerPeak float32 `json:"power_peak"`
	PowerPeakToday float32 `json:"power_peak_today"`
}

type Inverter struct {
	Db                   *gorm.DB                    `inject:""`
	SchedulerManager     *scheduler.SchedulerManager `inject:""`
	ConfigurationService *Configuration              `inject:""`
	NotificationService  *Notification               `inject:""`
	InfluxdbClient  *influxdb.Client                  `inject:""`
	Lock sync.Mutex
}

func (i *Inverter) Init() {
	i.SchedulerManager.ScheduleExecution(uint64(viper.GetInt("inverter.intervalSeconds")), i.ScheduledMeasurement)
}

func (i *Inverter) GetInverters(configurationID string) []model.Inverter {

	var Inverters []model.Inverter
	i.Db.Where("configuration_id = ?", configurationID).Find(&Inverters)

	return Inverters
}

func (i *Inverter) CreateOrUpdateInverter(configurationID string, Inverter *model.Inverter) {

	id, err := strconv.ParseUint(configurationID, 10, 32)
	if err == nil {
		Inverter.ConfigurationID = uint(id)
		i.Db.Save(&Inverter)
	}
}

func (i *Inverter) GetInverter(InverterID string) (*model.Inverter, error) {

	var Inverter model.Inverter
	i.Db.First(&Inverter, InverterID)
	if Inverter.ID == 0 {
		return nil, errors.New("Can't find Inverter with ID " + string(InverterID))
	}
	return &Inverter, nil
}

func (i *Inverter) DeleteInverter(InverterID string) error {

	var Inverter model.Inverter
	i.Db.First(&Inverter, InverterID)
	if Inverter.ID == 0 {
		return errors.New("Can't find Inverter with ID " + string(InverterID))
	}
	i.Db.Unscoped().Delete(&Inverter)
	return nil
}

func (i *Inverter) GetMetrics(Inverter *model.Inverter) (*InverterMetrics, error) {
	// Avoid accessing simultaneously inverter API
	result := &InverterMetrics{}
	i.Lock.Lock()
	resp, err := resty.R().SetResult(&result).Get(getInverterUrl(Inverter))
	i.Lock.Unlock()
	if err == nil && resp.StatusCode() == 200 {
		return result, nil
	} else if err == nil && resp.StatusCode() == 500 {
		return nil, errors.New("inverter is not available, probably powered off")
	} else {
		return nil, errors.New("Unable to get inverter metrics: " + err.Error())
	}

}

func getInverterUrl(inverter *model.Inverter) string {

	host := inverter.Host
	port := inverter.Port
	return fmt.Sprintf("http://%s:%d/metrics", host, port)
}

func (i *Inverter) ScheduledMeasurement() {

	configuration := i.ConfigurationService.GetCurrent()

	for _, inverter := range configuration.Inverters {

		log.Info("Getting inverter " + inverter.Name + " metrics")
		result, err := i.GetMetrics(&inverter)
		if err == nil {
			// Send data to influxdb
			point := client.Point{
				Measurement: "inverter",
				Tags: map[string]string{
					"name": inverter.Name,
				},
				Fields: map[string]interface{}{
					"power_pin_1": result.PowerPin1,
					"power_pin_2": result.PowerPin2,
					"grid_power_reading": result.GridPowerReading,
					"riso": result.Riso,
					"inverter_temperature": result.InverterTemperature,
					"booster_temperature": result.BoosterTemperature,
					"dc_ac_conversion_efficiency": result.DcAcConversionEfficiency,
				},
				Time:      time.Now(),
			}
			i.InfluxdbClient.AddPoint(point)
		} else {
			if !strings.HasPrefix(err.Error(), "inverter is not available") {
				log.Error("Unable to fetch inverter measurement. Reason:", err)
			}
		}
	}

}