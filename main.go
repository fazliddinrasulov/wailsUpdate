package main

import (
	"context"
	"embed"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)
//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	updater := NewUpdater("fazliddinrasulov/wailsUpdate")

	err := wails.Run(&options.App{
		Title:  "Your App",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			updater.Startup(ctx)
			
			// Check for updates 5 seconds after startup
			go func() {
				time.Sleep(5 * time.Second)
				updater.AutoCheckForUpdates()
			}()
		},
		Bind: []interface{}{
			app,
			updater,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}