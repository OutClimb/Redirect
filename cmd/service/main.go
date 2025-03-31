package main

import (
	"github.com/OutClimb/Redirect/internal/app"
	"github.com/OutClimb/Redirect/internal/http"
	"github.com/OutClimb/Redirect/internal/store"
)

func main() {
	storeLayer := store.New()
	appLayer := app.New(storeLayer)
	httpLayer := http.New(appLayer)

	httpLayer.Run()
}
