package app

import (
	"time"

	"github.com/OutClimb/Redirect/internal/store"
)

type RedirectInternal struct {
	ID       uint
	FromPath string
	ToUrl    string
	StartsOn *time.Time
	StopsOn  *time.Time
}

func (r *RedirectInternal) Internalize(redirect *store.Redirect) {
	r.ID = redirect.ID
	r.FromPath = redirect.FromPath
	r.ToUrl = redirect.ToUrl

	r.StartsOn = nil
	if redirect.StartsOn.Valid {
		r.StartsOn = &redirect.StartsOn.Time
	}

	r.StopsOn = nil
	if redirect.StopsOn.Valid {
		r.StopsOn = &redirect.StopsOn.Time
	}
}

func (a *appLayer) CreateRedirect(fromPath, toUrl string, startsOn, stopsOn int64) (*RedirectInternal, error) {
	startsOnTime := time.UnixMilli(startsOn)
	stopsOnTime := time.UnixMilli(stopsOn)

	if redirect, err := a.store.CreateRedirect(fromPath, toUrl, &startsOnTime, &stopsOnTime); err != nil {
		return &RedirectInternal{}, err
	} else {
		redirectInternal := RedirectInternal{}
		redirectInternal.Internalize(redirect)

		return &redirectInternal, nil
	}
}

func (a *appLayer) DeleteRedirect(id uint) error {
	if err := a.store.DeleteRedirect(id); err != nil {
		return err
	}

	return nil
}

func (a *appLayer) FindRedirect(path string) (*RedirectInternal, error) {
	if redirect, err := a.store.FindActiveRedirectByPath(path); err != nil {
		return &RedirectInternal{}, err
	} else {
		redirectInternal := RedirectInternal{}
		redirectInternal.Internalize(redirect)

		return &redirectInternal, nil
	}
}

func (a *appLayer) GetRedirect(id uint) (*RedirectInternal, error) {
	if redirect, err := a.store.GetRedirect(id); err != nil {
		return &RedirectInternal{}, err
	} else {
		redirectInternal := RedirectInternal{}
		redirectInternal.Internalize(redirect)

		return &redirectInternal, nil
	}
}

func (a *appLayer) GetAllRedirects() (*[]RedirectInternal, error) {
	if redirects, err := a.store.GetAllRedirects(); err != nil {
		return &[]RedirectInternal{}, err
	} else {
		redirectsInternal := make([]RedirectInternal, len(*redirects))
		for i, redirect := range *redirects {
			redirectsInternal[i].Internalize(&redirect)
		}

		return &redirectsInternal, nil
	}
}

func (a *appLayer) UpdateRedirect(id uint, fromPath, toUrl string, startsOn, stopsOn int64) (*RedirectInternal, error) {
	startsOnTime := time.UnixMilli(startsOn)
	stopsOnTime := time.UnixMilli(stopsOn)

	if redirect, err := a.store.UpdateRedirect(id, fromPath, toUrl, &startsOnTime, &stopsOnTime); err != nil {
		return &RedirectInternal{}, err
	} else {
		redirectInternal := RedirectInternal{}
		redirectInternal.Internalize(redirect)

		return &redirectInternal, nil
	}
}
