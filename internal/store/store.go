package store

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type StoreLayer interface {
	// Redirect
	CreateRedirect(fromPath, toUrl string, startsOn, stopsOn *time.Time) (*Redirect, error)
	DeleteRedirect(id uint) error
	FindActiveRedirectByPath(path string) (*Redirect, error)
	GetAllRedirects() (*[]Redirect, error)
	GetRedirect(id uint) (*Redirect, error)
	UpdateRedirect(id uint, fromPath, toUrl string, startsOn, stopsOn *time.Time) (*Redirect, error)

	// User
	CreateUser(username, password, name, email, role string) (*User, error)
	GetUser(id uint) (*User, error)
	GetUserWithUsername(username string) (*User, error)
	UpdatePassword(id uint, password string) error
}

type storeLayer struct {
	db *gorm.DB
}

func New() *storeLayer {
	username := os.Getenv("DB_USERNAME")
	if len(username) == 0 {
		log.Fatal("Error: No database username provided")
		return nil
	}

	password := os.Getenv("DB_PASSWORD")

	host := os.Getenv("DB_HOST")
	if len(host) == 0 {
		log.Fatal("Error: No database hostname provided")
		return nil
	}

	name := os.Getenv("DB_NAME")
	if len(name) == 0 {
		log.Fatal("Error: No database name provided")
		return nil
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, name)), &gorm.Config{})
	if err != nil {
		log.Fatal("Error: Unable to connect to MySQL server", err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Redirect{})

	return &storeLayer{
		db: db,
	}
}
