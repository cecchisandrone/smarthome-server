package service

import (
	"errors"
	"fmt"
	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/jinzhu/gorm"
	"gopkg.in/resty.v1"
	"strconv"
)

type WellPump struct {
	Db *gorm.DB `inject:""`
}

func (w *WellPump) Init() {
}

func (w *WellPump) GetWellPumps(configurationID string) []model.WellPump {

	var WellPumps []model.WellPump
	w.Db.Where("configuration_id = ?", configurationID).Find(&WellPumps)

	return WellPumps
}

func (w *WellPump) CreateOrUpdateWellPump(configurationID string, wellPump *model.WellPump) {

	id, err := strconv.ParseUint(configurationID, 10, 32)
	if err == nil {
		wellPump.ConfigurationID = uint(id)
		w.Db.Save(&wellPump)
	}
}

func (w *WellPump) GetWellPump(wellPumpID string) (*model.WellPump, error) {

	var wellPump model.WellPump
	w.Db.First(&wellPump, wellPumpID)
	if wellPump.ID == 0 {
		return nil, errors.New("Can't find WellPump with ID " + string(wellPumpID))
	}
	return &wellPump, nil
}

func (w *WellPump) DeleteWellPump(wellPumpID string) error {

	var wellPump model.WellPump
	w.Db.First(&wellPump, wellPumpID)
	if wellPump.ID == 0 {
		return errors.New("Can't find WellPump with ID " + string(wellPumpID))
	}
	w.Db.Unscoped().Delete(&wellPump)
	return nil
}

func (w *WellPump) GetRelay(wellPump *model.WellPump) (int, error) {
	result := make(map[string]int)
	resp, err := resty.R().SetResult(&result).Get(getWellPumpUrl(wellPump))
	if err == nil && resp.StatusCode() == 200 {
		return result["status"], nil
	} else {
		return -1, errors.New("unable to get relay status")
	}
}

func (w *WellPump) ToggleRelay(wellPump *model.WellPump, status int) error {
	body := map[string]interface{}{"status": status}
	resp, err := resty.R().SetBody(body).Put(getWellPumpUrl(wellPump))
	if err != nil || resp.StatusCode() != 200 {
		return errors.New(fmt.Sprintf("unable to toggle well pump %s to %d", wellPump.Name, status))
	}
	return nil
}

func getWellPumpUrl(wellPump *model.WellPump) string {

	host := wellPump.Host
	port := wellPump.Port
	return fmt.Sprintf("http://%s:%d/relay", host, port)
}
