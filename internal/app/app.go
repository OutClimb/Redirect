package app

import (
	"github.com/OutClimb/Redirect/internal/store"
)

type AppLayer interface {
	// Redirect
	CreateRedirect(fromPath, toUrl string, startsOn, stopsOn int64) (*RedirectInternal, error)
	DeleteRedirect(id uint) error
	FindRedirect(path string) (*RedirectInternal, error)
	GetRedirect(id uint) (*RedirectInternal, error)
	GetAllRedirects() (*[]RedirectInternal, error)
	UpdateRedirect(id uint, fromPath, toUrl string, startsOn, stopsOn int64) (*RedirectInternal, error)

	// User
	AuthenticateUser(username, password string) (*UserInternal, error)
	CheckRole(userRole, requiredRole string) bool
	CreateToken(user *UserInternal, clientIp string) (string, error)
	GetUser(userId uint) (*UserInternal, error)
	ValidatePassword(user *UserInternal, password string) error
	ValidateUser(userId uint) error
	UpdatePassword(user *UserInternal, password string) error
}

type appLayer struct {
	store store.StoreLayer
}

func New(storeLayer store.StoreLayer) *appLayer {
	return &appLayer{
		store: storeLayer,
	}
}
