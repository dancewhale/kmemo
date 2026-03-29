package main

import (
	"context"

	"kmemo/internal/app"
	"kmemo/internal/bootstrap"
)

// App is the Wails binding surface. The name "App" maps to window.go.main.App in the frontend.
type App struct {
	*app.Desktop
}

func newApp() (*App, error) {
	d, err := bootstrap.NewDesktop(context.Background())
	if err != nil {
		return nil, err
	}
	return &App{Desktop: d}, nil
}
