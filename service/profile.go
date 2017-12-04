package service

import (
	"errors"

	"fmt"

	"github.com/cecchisandrone/smarthome-server/model"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Profile struct {
	Db *gorm.DB `inject:""`
}

func (p Profile) Init() {

}

func (p Profile) GetProfiles() ([]model.Profile, error) {

	var profiles []model.Profile
	err := p.Db.Find(&profiles).Error
	return profiles, err
}

func (p Profile) CreateProfile(profile *model.Profile) error {

	if p.Db.Where("username = ?", profile.Username).First(&model.Profile{}).RecordNotFound() {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(profile.Password), bcrypt.DefaultCost)
		profile.Password = string(hashedPassword)
		return p.Db.Save(profile).Error
	}
	return errors.New("User already exists")
}

func (p Profile) GetProfile(profileID string) (*model.Profile, error) {

	var profile model.Profile
	if p.Db.First(&profile, profileID).RecordNotFound() {
		return nil, errors.New("Can't find profile with ID " + string(profileID))
	}
	return &profile, nil
}

func (p Profile) Authenticate(username string, password string) bool {
	var profile model.Profile
	if !p.Db.Where("username = ?", username).First(&profile).RecordNotFound() {
		if err := bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(password)); err == nil {
			fmt.Println(err)
			return true
		}
	}
	return false
}

func (p Profile) GetProfileByUsername(username string) (*model.Profile, error) {
	var profile model.Profile
	if !p.Db.Where("username = ?", username).First(&profile).RecordNotFound() {
		return &profile, nil
	}
	return nil, p.Db.Error
}
