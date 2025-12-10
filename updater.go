package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Updater handles GitHub-based auto-updates
type Updater struct {
	ctx        context.Context
	repoOwner  string
	repoName   string
	currentVer string
}

// GitHubRelease represents a GitHub release API response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
	PublishedAt string `json:"published_at"`
}

// UpdateResult is returned to the frontend
type UpdateResult struct {
	Available   bool   `json:"available"`
	Version     string `json:"version"`
	ReleaseDate string `json:"release_date"`
	Changelog   string `json:"changelog"`
	DownloadURL string `json:"download_url"`
}

// NewUpdater creates a new Updater instance
// repo should be in format "owner/repo"
func NewUpdater(repo string) *Updater {
	return &Updater{
		repoOwner:  repo,
		currentVer: "1.0.0", // Change this to your app version
	}
}

// Startup is called when the app starts
func (u *Updater) Startup(ctx context.Context) {
	u.ctx = ctx
}

// AutoCheckForUpdates checks for updates silently and shows dialog if available
func (u *Updater) AutoCheckForUpdates() {
	result, err := u.CheckForUpdates()
	if err != nil {
		return // Silent fail
	}
	if result.Available {
		// Show dialog to user
		wailsRuntime.EventsEmit(u.ctx, "update-available", result)
	}
}

// CheckForUpdates checks GitHub for the latest release
func (u *Updater) CheckForUpdates() (*UpdateResult, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", u.repoOwner)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse release info: %w", err)
	}

	// Find appropriate asset for current OS
	var downloadURL string
	osName := runtime.GOOS
	for _, asset := range release.Assets {
		if osName == "windows" && (contains(asset.Name, "windows") || contains(asset.Name, ".exe")) {
			downloadURL = asset.BrowserDownloadURL
			break
		} else if osName == "darwin" && (contains(asset.Name, "darwin") || contains(asset.Name, "mac") || contains(asset.Name, ".dmg")) {
			downloadURL = asset.BrowserDownloadURL
			break
		} else if osName == "linux" && contains(asset.Name, "linux") {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	result := &UpdateResult{
		Available:   release.TagName != u.currentVer && release.TagName != "v"+u.currentVer,
		Version:     release.TagName,
		ReleaseDate: release.PublishedAt,
		Changelog:   release.Body,
		DownloadURL: downloadURL,
	}

	return result, nil
}

// DownloadUpdate downloads the update file
func (u *Updater) DownloadUpdate(url, destPath string) error {
	client := &http.Client{Timeout: 5 * time.Minute}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status: %d", resp.StatusCode)
	}

	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// InstallUpdate runs the installer
func (u *Updater) InstallUpdate(installerPath string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command(installerPath, "/SILENT")
		return cmd.Start()
	case "darwin":
		cmd := exec.Command("open", installerPath)
		return cmd.Start()
	case "linux":
		os.Chmod(installerPath, 0755)
		cmd := exec.Command(installerPath)
		return cmd.Start()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
