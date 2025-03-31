package store

import (
	"database/sql"
	"time"
)

type Redirect struct {
	ID       uint   `gorm:"primaryKey"`
	FromPath string `gorm:"not null"`
	ToUrl    string `gorm:"not null"`
	StartsOn sql.NullTime
	StopsOn  sql.NullTime
}

func (s *storeLayer) CreateRedirect(fromPath, toUrl string, startsOn, stopsOn *time.Time) (*Redirect, error) {
	redirect := Redirect{
		FromPath: fromPath,
		ToUrl:    toUrl,
		StartsOn: sql.NullTime{Time: *startsOn, Valid: startsOn != nil && startsOn.UnixMilli() != 0},
		StopsOn:  sql.NullTime{Time: *stopsOn, Valid: stopsOn != nil && startsOn.UnixMilli() != 0},
	}

	if result := s.db.Create(&redirect); result.Error != nil {
		return nil, result.Error
	}

	return &redirect, nil
}

func (s *storeLayer) DeleteRedirect(id uint) error {
	if result := s.db.Delete(&Redirect{}, id); result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *storeLayer) FindActiveRedirectByPath(path string) (*Redirect, error) {
	redirect := Redirect{}

	if result := s.db.Model(&Redirect{}).Where("from_path = ? AND (starts_on IS NULL or starts_on <= NOW()) AND (stops_on IS NULL or stops_on > NOW())", path).First(&redirect); result.Error != nil {
		return &Redirect{}, result.Error
	}

	return &redirect, nil
}

func (s *storeLayer) GetAllRedirects() (*[]Redirect, error) {
	redirects := []Redirect{}

	if result := s.db.Find(&redirects); result.Error != nil {
		return &[]Redirect{}, result.Error
	}

	return &redirects, nil
}

func (s *storeLayer) GetRedirect(id uint) (*Redirect, error) {
	redirect := Redirect{}

	if result := s.db.Where("id = ?", id).First(&redirect); result.Error != nil {
		return &Redirect{}, result.Error
	}

	return &redirect, nil
}

func (s *storeLayer) UpdateRedirect(id uint, fromPath, toUrl string, startsOn, stopsOn *time.Time) (*Redirect, error) {
	redirect, err := s.GetRedirect(id)
	if err != nil {
		return nil, err
	}

	redirect.FromPath = fromPath
	redirect.ToUrl = toUrl
	redirect.StartsOn = sql.NullTime{Time: *startsOn, Valid: startsOn != nil && startsOn.UnixMilli() != 0}
	redirect.StopsOn = sql.NullTime{Time: *stopsOn, Valid: stopsOn != nil && startsOn.UnixMilli() != 0}

	if result := s.db.Save(&redirect); result.Error != nil {
		return nil, result.Error
	}

	return redirect, nil
}
