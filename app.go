package main

import (
	"context"
	"fmt"
	"os/exec"
	"net/http"
	"encoding/json"
	"runtime"
	"os"
	"io"
	"path/filepath"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SilentPrint prints a document silently based on the OS
func (a *App) SilentPrint(path string) error {
	switch runtime.GOOS {
	case "windows":
		return a.silentPrintWindows(path)
	case "darwin":
		return a.silentPrintMac(path)
	case "linux":
		return a.silentPrintLinux(path)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// silentPrintWindows prints on Windows using PowerShell
func (a *App) silentPrintWindows(path string) error {
	cmd := exec.Command("powershell", "-command", "Start-Process", path, "-Verb", "Print", "-WindowStyle", "Hidden")
	return cmd.Run()
}

// silentPrintMac prints on macOS using lp command
func (a *App) silentPrintMac(path string) error {
	cmd := exec.Command("lp", path)
	return cmd.Run()
}

// silentPrintLinux prints on Linux using lp command
func (a *App) silentPrintLinux(path string) error {
	cmd := exec.Command("lp", path)
	return cmd.Run()
}

// UpdateInfo contains information about available updates
type UpdateInfo struct {
	Version     string `json:"version"`
	URL         string `json:"url"`
	ReleaseDate string `json:"release_date"`
	ChangeLog   string `json:"changelog"`
}

// CheckUpdate checks for available updates from GitHub
func (a *App) CheckUpdate(currentVersion string) (*UpdateInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://raw.githubusercontent.com/fazliddinrasulov/wailsUpdate/main/main/latest.json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch update info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var info UpdateInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode update info: %w", err)
	}

	if info.Version != currentVersion {
		return &info, nil
	}

	return nil, nil
}

// DownloadUpdate downloads the update file to the specified path
func (a *App) DownloadUpdate(url, destPath string) error {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create destination directory if it doesn't exist
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// InstallUpdate installs the downloaded update
func (a *App) InstallUpdate(installerPath string) error {
	switch runtime.GOOS {
	case "windows":
		// Execute the installer
		cmd := exec.Command(installerPath, "/SILENT")
		return cmd.Start()
	case "darwin":
		// Open the DMG or execute installer
		cmd := exec.Command("open", installerPath)
		return cmd.Start()
	case "linux":
		// Make executable and run
		os.Chmod(installerPath, 0755)
		cmd := exec.Command(installerPath)
		return cmd.Start()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// GetAppVersion returns the current application version
func (a *App) GetAppVersion() string {
	return "1.0.0" // Update this with your actual version
}