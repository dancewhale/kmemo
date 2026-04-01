package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"go.uber.org/zap"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app, err := newApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "kmemo desktop bootstrap failed: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if app.logger != nil {
			_ = app.logger.Sync()
		}
	}()

	err = wails.Run(&options.App{
		Title:  "kmemo",
		Width:  1024,
		Height: 640,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.OnStartup,
		OnShutdown:       app.OnShutdown,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{},
		Mac:     &mac.Options{},
		Linux:   &linux.Options{},
	})
	if err != nil {
		app.logger.Error("wails run failed", zap.Error(err))
		os.Exit(1)
	}
}
