package service

import (
	"errors"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/jinzhu/gorm"
)

type Profile struct {
	Db *gorm.DB `inject:""`
}

func (p Profile) GetProfiles() []model.Profile {

	var profiles []model.Profile
	p.Db.Find(&profiles)
	return profiles
}

func (p Profile) CreateProfile(profile *model.Profile) {

	p.Db.Save(&profile)
}

func (p Profile) GetProfile(profileID string) (*model.Profile, error) {

	var profile model.Profile
	p.Db.First(&profile, profileID)
	if profile.ID == 0 {
		return nil, errors.New("Can't find profile with ID " + string(profileID))
	}
	return &profile, nil
}
