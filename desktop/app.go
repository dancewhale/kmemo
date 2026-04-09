package main

import (
	"context"

	"go.uber.org/zap"

	"kmemo/internal/app"
	"kmemo/internal/bootstrap"
)

// App is the Wails binding surface. The name "App" maps to window.go.main.App in the frontend.
type App struct {
	*app.Desktop
	logger *zap.Logger
}

func newApp() (*App, error) {
	h, err := bootstrap.NewHeadless(context.Background())
	if err != nil {
		return nil, err
	}
	return &App{
		Desktop: app.NewDesktop(h.Config, h.Logger, h.Worker),
		logger:  h.Logger,
	}, nil
}
