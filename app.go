package main

import (
	"FanslyArchivePlayer/backend/integrations"
	"FanslyArchivePlayer/backend/integrations/fansly"
	"FanslyArchivePlayer/backend/models"
	"FanslyArchivePlayer/backend/services"
	"context"
	_ "database/sql"
	_ "encoding/json"
	_ "errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/pelletier/go-toml"
	"github.com/sqweek/dialog"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"net/http"
	"net/url"
	"os"
	_ "os/exec"
	"path/filepath"
	_ "runtime"
	_ "strconv"
	"strings"
)

// App struct
type App struct {
	ctx               context.Context
	videoService      *services.VideoService
	fileDialogService *services.FileDialogService
	cacheService      *services.CacheService
	integrations      *integrations.Manager
	currentVideoPath  string
	appDataDir        string
}

// BrowseForFile opens a file dialog for selecting a file
func (a *App) BrowseForFile(title string, filter string) (string, error) {
	// Parse filter string (format: "Description|*.ext")
	filterParts := strings.Split(filter, "|")
	if len(filterParts) != 2 {
		filter = ""
	}
	// Set dialog title
	dialog.Directory().Title(title)
	// Open file dialog
	if filter != "" {
		return dialog.File().Filter(filterParts[0], filterParts[1]).Load()
	}
	return dialog.File().Load()
}

// BrowseForFolder opens a folder dialog for selecting a directory
func (a *App) BrowseForFolder(title string) (string, error) {
	// Set dialog title
	dialog.Directory().Title(title)
	// Open folder dialog
	return dialog.Directory().Browse()
}

// NewApp creates a new App application struct
func NewApp() *App {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	appDataDir := filepath.Join(configDir, "ArchivePlayer")
	// Create app data directory if it doesn't exist
	if _, err := os.Stat(appDataDir); os.IsNotExist(err) {
		err = os.MkdirAll(appDataDir, 0755)
		if err != nil {
			// Fall back to current directory if we can't create the app data directory
			appDataDir = "."
		}
	}

	cacheService := services.NewCacheService(appDataDir)

	return &App{
		videoService:      services.NewVideoService(),
		fileDialogService: services.NewFileDialogService(),
		cacheService:      cacheService,
		integrations:      integrations.NewManager(appDataDir, cacheService),
		appDataDir:        appDataDir,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.fileDialogService.SetContext(ctx)
	// Set up a handler for serving video files
	http.HandleFunc("/video/", func(w http.ResponseWriter, r *http.Request) {
		if a.currentVideoPath == "" {
			http.Error(w, "No video loaded", http.StatusNotFound)
			return
		}
		// Set appropriate headers
		w.Header().Set("Content-Type", getContentType(a.currentVideoPath))
		http.ServeFile(w, r, a.currentVideoPath)
	})
	// Set up a handler for serving thumbnail images
	http.HandleFunc("/thumbnail/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the file path from the URL
		urlPath := r.URL.Path
		if !strings.HasPrefix(urlPath, "/thumbnail/") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		// Decode the file path
		filePath, err := decodeFilePath(strings.TrimPrefix(urlPath, "/thumbnail/"))
		if err != nil {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}
		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "Thumbnail not found", http.StatusNotFound)
			return
		}
		// Set appropriate headers
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, filePath)
	})
	// Start the HTTP server
	go http.ListenAndServe(":8080", nil)
}

// GetFanslyStreams retrieves all streams from the Fansly database
func (a *App) GetFanslyStreams() (fansly.StreamsResult, error) {
	return a.integrations.FanslyService.GetStreams()
}

// LoadFanslyStream loads a Fansly stream and its associated chat file
func (a *App) LoadFanslyStream(streamPath string) (fansly.StreamResult, error) {
	return a.integrations.FanslyService.LoadStream(streamPath)
}

// GetFanslyConfig returns the current Fansly integration configuration
func (a *App) GetFanslyConfig() (fansly.Config, error) {
	return a.integrations.FanslyService.GetConfig()
}

// SaveFanslyConfig saves the Fansly integration configuration
func (a *App) SaveFanslyConfig(config fansly.Config) error {
	return a.integrations.FanslyService.SaveConfig(config)
}

// OpenVideoFile opens a dialog to select a video file
func (a *App) OpenVideoFile() (string, error) {
	path, err := a.fileDialogService.OpenVideoFile()
	if err != nil {
		return "", err
	}
	return a.LoadVideoFromPath(path)
}

// LoadVideoFromPath loads a video from a specific path
func (a *App) LoadVideoFromPath(path string) (string, error) {
	err := a.videoService.LoadVideo(path)
	if err != nil {
		return "", err
	}
	// Store the current video path
	a.currentVideoPath = path
	// Return a URL that can be used by the video element
	return "http://localhost:8080/video/" + filepath.Base(path), nil
}

// OpenChatFile opens a dialog to select a chat JSON file
func (a *App) OpenChatFile(path []string) (string, error) {
	var chatFilePath string
	var err error
	if len(path) > 0 && path[0] != "" {
		// Use the provided path
		chatFilePath = path[0]
	} else {
		// Open file dialog
		chatFilePath, err = wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
			Title: "Select Chat File",
			Filters: []wailsRuntime.FileFilter{
				{
					DisplayName: "JSON Files (*.json)",
					Pattern:     "*.json",
				},
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to open file dialog: %v", err)
		}
	}
	if chatFilePath == "" {
		return "", nil // User cancelled
	}
	// Load the chat file
	err = a.videoService.LoadChatFile(chatFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to load chat file: %v", err)
	}
	// Return the URL for the chat file
	return fmt.Sprintf("http://localhost:8080/video/%s", filepath.Base(chatFilePath)), nil
}

// LoadChatFromPath loads a chat file from a specific path
func (a *App) LoadChatFromPath(path string) (string, error) {
	err := a.videoService.LoadChatFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load chat file: %v", err)
	}
	return path, nil
}

// GetMessagesAtTime returns messages at a specific time
func (a *App) GetMessagesAtTime(currentTime float64, windowSize float64) []models.ChatMessage {
	return a.videoService.GetMessagesAtTime(currentTime, windowSize)
}

// GetVideoFileInfo returns information about the current video
func (a *App) GetVideoFileInfo() map[string]string {
	return a.videoService.GetVideoFileInfo()
}

// GetAllChatMessages returns all chat messages
func (a *App) GetAllChatMessages() []models.ChatMessage {
	return a.videoService.ChatMessages
}

// Helper function to get content type from file extension
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".ogg":
		return "video/ogg"
	case ".mkv":
		return "video/x-matroska"
	default:
		return "application/octet-stream"
	}
}

// Helper function to decode file path from URL
func decodeFilePath(encodedPath string) (string, error) {
	// URL decode the path
	decodedPath, err := url.QueryUnescape(encodedPath)
	if err != nil {
		return "", err
	}
	return decodedPath, nil
}
