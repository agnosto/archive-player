package fansly

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	_ "time"

	"FanslyArchivePlayer/backend/services"
	//_ "github.com/mattn/go-sqlite3"
	"github.com/pelletier/go-toml"
	_ "modernc.org/sqlite"
)

// Service handles Fansly integration functionality
type Service struct {
	appDataDir   string
	cacheService *services.CacheService
}

// NewService creates a new Fansly service
func NewService(appDataDir string, cacheService *services.CacheService) *Service {
	return &Service{
		appDataDir:   appDataDir,
		cacheService: cacheService,
	}
}

// GetConfig returns the current Fansly integration configuration
func (s *Service) GetConfig() (Config, error) {
	config := Config{}
	// Try to load from app settings
	configBytes, err := os.ReadFile(filepath.Join(s.appDataDir, "fansly_config.json"))
	if err == nil {
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			return config, err
		}
		return config, nil
	}

	// If not found, try to detect default locations
	defaultConfigPath := getDefaultConfigPath()
	if _, err := os.Stat(defaultConfigPath); err == nil {
		config.ConfigPath = defaultConfigPath
		// Try to parse the config to find the save location
		if configContent, err := os.ReadFile(defaultConfigPath); err == nil {
			parsedConfig, err := toml.Load(string(configContent))
			if err == nil {
				if options := parsedConfig.Get("options"); options != nil {
					if saveLocation, ok := options.(*toml.Tree).Get("save_location").(string); ok && saveLocation != "" {
						dbPath := filepath.Join(saveLocation, "downloads.db")
						if _, err := os.Stat(dbPath); err == nil {
							config.DbPath = saveLocation
						}
					}
				}
			}
		}
	}
	return config, nil
}

// SaveConfig saves the Fansly integration configuration
func (s *Service) SaveConfig(config Config) error {
	// Validate config
	if config.ConfigPath == "" {
		return errors.New("config path cannot be empty")
	}
	if _, err := os.Stat(config.ConfigPath); err != nil {
		return errors.New("config file not found at specified path")
	}

	// If dbPath is provided, check if downloads.db exists
	if config.DbPath != "" {
		dbFilePath := filepath.Join(config.DbPath, "downloads.db")
		if _, err := os.Stat(dbFilePath); err != nil {
			return errors.New("downloads.db not found in the specified folder")
		}
	} else {
		// Try to extract dbPath from config file
		if configContent, err := os.ReadFile(config.ConfigPath); err == nil {
			parsedConfig, err := toml.Load(string(configContent))
			if err == nil {
				if options := parsedConfig.Get("options"); options != nil {
					if saveLocation, ok := options.(*toml.Tree).Get("save_location").(string); ok && saveLocation != "" {
						dbPath := filepath.Join(saveLocation, "downloads.db")
						if _, err := os.Stat(dbPath); err == nil {
							config.DbPath = saveLocation
						} else {
							return errors.New("could not find downloads.db in the save_location from config")
						}
					}
				}
			}
		}
		if config.DbPath == "" {
			return errors.New("could not determine database path")
		}
	}

	// Save config to app data directory
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.appDataDir, "fansly_config.json"), configBytes, 0644)
}

// GetStreams retrieves all streams from the Fansly database
func (s *Service) GetStreams() (StreamsResult, error) {
	result := StreamsResult{
		Streams:   []Stream{},
		ChatFiles: []string{},
	}

	// Get Fansly config
	config, err := s.GetConfig()
	if err != nil {
		result.Error = "Failed to load Fansly configuration"
		return result, nil
	}
	if config.DbPath == "" {
		result.Error = "Fansly integration not configured. Please set up the integration first."
		return result, nil
	}

	// Load video cache
	videoCache, err := s.cacheService.LoadVideoCache("fansly")
	if err != nil {
		// Log error but continue without cache
		fmt.Printf("Failed to load video cache: %v\n", err)
		videoCache = services.VideoCache{
			Videos: make(map[string]services.VideoMetadata),
		}
	}

	dbPath := filepath.Join(config.DbPath, "downloads.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		result.Error = "Failed to open downloads database: " + err.Error()
		return result, nil
	}
	defer db.Close()

	// Query for livestreams and contact sheets
	rows, err := db.Query("SELECT model, hash, path, file_type FROM files WHERE file_type IN ('livestream', 'contact_sheet')")
	if err != nil {
		result.Error = "Failed to query database: " + err.Error()
		return result, nil
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var stream Stream
		err := rows.Scan(&stream.Model, &stream.Hash, &stream.Path, &stream.FileType)
		if err != nil {
			continue
		}

		// Verify file exists
		fileInfo, err := os.Stat(stream.Path)
		if err != nil {
			continue
		}

		result.Streams = append(result.Streams, stream)

		// If this is a livestream, check if we need to update the cache
		if stream.FileType == "livestream" {
			// Check if we have this video in cache and if it's still valid
			cachedVideo, exists := videoCache.Videos[stream.Path]
			needsUpdate := !exists ||
				cachedVideo.LastModified != fileInfo.ModTime() ||
				cachedVideo.FileSize != fileInfo.Size() ||
				cachedVideo.Duration == 0

			if needsUpdate {
				// Get video duration if possible
				if duration, err := getVideoDuration(stream.Path); err == nil {
					// Update cache
					videoCache.Videos[stream.Path] = services.VideoMetadata{
						Path:         stream.Path,
						Hash:         stream.Hash,
						Duration:     duration,
						LastModified: fileInfo.ModTime(),
						FileSize:     fileInfo.Size(),
					}
				}
			}
		}
	}

	// Save updated cache
	if err := s.cacheService.SaveVideoCache("fansly", videoCache); err != nil {
		fmt.Printf("Failed to save video cache: %v\n", err)
	}

	// Find chat files and add duration from cache
	for i, stream := range result.Streams {
		if stream.FileType == "livestream" {
			// Check for chat file
			chatPath := strings.TrimSuffix(stream.Path, filepath.Ext(stream.Path)) + "_chat.json"
			if _, err := os.Stat(chatPath); err == nil {
				result.ChatFiles = append(result.ChatFiles, chatPath)
			}

			// Check for contact sheet
			contactSheetPath := strings.TrimSuffix(stream.Path, filepath.Ext(stream.Path)) + "_contact_sheet.jpg"
			if _, err := os.Stat(contactSheetPath); err == nil {
				// Update the stream in the result with the contact sheet path
				result.Streams[i].ContactSheet = contactSheetPath
			}

			// Get duration from cache
			if cachedVideo, exists := videoCache.Videos[stream.Path]; exists && cachedVideo.Duration > 0 {
				result.Streams[i].Duration = cachedVideo.Duration
			}
		}
	}

	return result, nil
}

// LoadStream loads a Fansly stream and its associated chat file
func (s *Service) LoadStream(streamPath string) (StreamResult, error) {
	result := StreamResult{
		VideoPath: streamPath,
		Success:   false,
	}

	// Verify video file exists
	if _, err := os.Stat(streamPath); err != nil {
		result.Error = "Video file not found"
		return result, nil
	}

	// Check for chat file
	chatPath := strings.TrimSuffix(streamPath, filepath.Ext(streamPath)) + "_chat.json"
	if _, err := os.Stat(chatPath); err == nil {
		result.ChatPath = chatPath
	}

	// Check for contact sheet
	contactSheetPath := strings.TrimSuffix(streamPath, filepath.Ext(streamPath)) + "_contact_sheet.jpg"
	if _, err := os.Stat(contactSheetPath); err == nil {
		result.ContactSheet = contactSheetPath
	}

	result.Success = true
	return result, nil
}

// Helper function to get the default Fansly config path
func getDefaultConfigPath() string {
	currentDirConfig := "config.toml"
	if _, err := os.Stat(currentDirConfig); err == nil {
		return currentDirConfig
	}

	var configDir string
	var err error
	if runtime.GOOS == "darwin" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		configDir = filepath.Join(homeDir, ".config")
	} else {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return ""
		}
	}
	return filepath.Join(configDir, "fansly-scraper", "config.toml")
}

// Helper function to get video duration
func getVideoDuration(videoPath string) (float64, error) {
	// Check if ffprobe is available
	_, err := exec.LookPath("ffprobe")
	if err != nil {
		return 0, errors.New("ffprobe not found")
	}

	// Run ffprobe to get duration
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// Parse duration
	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}
