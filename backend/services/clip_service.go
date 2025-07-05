package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ClipStorageOption defines where clips should be stored
type ClipStorageOption string

const (
	// StoreInVideosDir saves clips in user's Videos/fansly-clips directory
	StoreInVideosDir ClipStorageOption = "videos_dir"
	// StoreWithSourceVideo saves clips in a clips folder next to the source video
	StoreWithSourceVideo ClipStorageOption = "source_video_dir"
	// StoreInCustomDir saves clips in a user-specified directory
	StoreInCustomDir ClipStorageOption = "custom_dir"
)

// ClipService handles video clipping functionality
type ClipService struct {
	appDataDir      string
	defaultOption   ClipStorageOption
	customOutputDir string
}

// NewClipService creates a new clip service
func NewClipService(appDataDir string) *ClipService {
	// Create default clips directory in app data as fallback
	clipsDir := filepath.Join(appDataDir, "clips")
	if _, err := os.Stat(clipsDir); os.IsNotExist(err) {
		os.MkdirAll(clipsDir, 0755)
	}

	return &ClipService{
		appDataDir:    appDataDir,
		defaultOption: StoreInVideosDir, // Default to user's Videos directory
	}
}

// SetStorageOption sets the preferred storage option for clips
func (s *ClipService) SetStorageOption(option ClipStorageOption, customDir string) {
	s.defaultOption = option
	if option == StoreInCustomDir && customDir != "" {
		s.customOutputDir = customDir
	}
}

// ClipResult contains information about the created clip
type ClipResult struct {
	Success      bool   `json:"success"`
	FilePath     string `json:"filePath"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// CreateClip creates a video clip from the source video
func (s *ClipService) CreateClip(sourceVideoPath string, startTime float64, duration float64, title string) ClipResult {
	// Validate inputs
	if sourceVideoPath == "" {
		return ClipResult{Success: false, ErrorMessage: "No source video provided"}
	}
	if duration <= 0 {
		return ClipResult{Success: false, ErrorMessage: "Duration must be a positive number"}
	}

	// Create a filename based on title or timestamp if title is empty
	filename := title
	if filename == "" {
		filename = fmt.Sprintf("clip_%s", time.Now().Format("20060102_150405"))
	}

	// Sanitize filename
	filename = sanitizeFilename(filename)

	// Determine output directory based on storage option
	outputDir, err := s.getOutputDirectory(sourceVideoPath)
	if err != nil {
		return ClipResult{Success: false, ErrorMessage: fmt.Sprintf("Failed to create output directory: %v", err)}
	}

	outputPath := filepath.Join(outputDir, filename+".mp4")

	// Format start time for ffmpeg (convert seconds to HH:MM:SS.mmm format)
	startTimeStr := formatFFmpegTime(startTime)
	durationStr := formatFFmpegTime(duration)

	// Create the clip using ffmpeg
	cmd := exec.Command(
		"ffmpeg",
		"-ss", startTimeStr,
		"-i", sourceVideoPath,
		"-t", durationStr,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-strict", "experimental",
		"-b:a", "128k",
		"-y", // Overwrite output file if it exists
		outputPath,
	)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ClipResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("FFmpeg error: %v\nOutput: %s", err, string(output)),
		}
	}

	return ClipResult{
		Success:  true,
		FilePath: outputPath,
	}
}

// getOutputDirectory determines where to save the clip based on the storage option
func (s *ClipService) getOutputDirectory(sourceVideoPath string) (string, error) {
	var outputDir string

	switch s.defaultOption {
	case StoreInVideosDir:
		// Save in user's Videos/fansly-clips directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not find user home directory: %v", err)
		}

		// Try to find the Videos directory
		videosDir := filepath.Join(homeDir, "Videos")
		if _, err := os.Stat(videosDir); os.IsNotExist(err) {
			// If Videos doesn't exist, try other common names
			videosDir = filepath.Join(homeDir, "videos")
			if _, err := os.Stat(videosDir); os.IsNotExist(err) {
				// If still not found, use Documents
				videosDir = filepath.Join(homeDir, "Documents")
				if _, err := os.Stat(videosDir); os.IsNotExist(err) {
					// Last resort, use home directory
					videosDir = homeDir
				}
			}
		}

		outputDir = filepath.Join(videosDir, "fansly-clips")

	case StoreWithSourceVideo:
		// Save in a clips folder next to the source video
		sourceDir := filepath.Dir(sourceVideoPath)
		outputDir = filepath.Join(sourceDir, "clips")

	case StoreInCustomDir:
		// Save in user-specified directory
		if s.customOutputDir != "" {
			outputDir = s.customOutputDir
		} else {
			// Fallback to app data directory if custom dir is not set
			outputDir = filepath.Join(s.appDataDir, "clips")
		}

	default:
		// Fallback to app data directory
		outputDir = filepath.Join(s.appDataDir, "clips")
	}

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create output directory: %v", err)
		}
	}

	return outputDir, nil
}

// GetClips returns a list of all saved clips
func (s *ClipService) GetClips() []string {
	var allClips []string

	// Get clips from all possible locations

	// 1. Check app data directory
	appDataClips := s.getClipsFromDir(filepath.Join(s.appDataDir, "clips"))
	allClips = append(allClips, appDataClips...)

	// 2. Check Videos/fansly-clips directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		videosDir := filepath.Join(homeDir, "Videos")
		if _, err := os.Stat(videosDir); err == nil {
			videosClips := s.getClipsFromDir(filepath.Join(videosDir, "fansly-clips"))
			allClips = append(allClips, videosClips...)
		}

		// Try lowercase "videos" too
		videosDir = filepath.Join(homeDir, "videos")
		if _, err := os.Stat(videosDir); err == nil {
			videosClips := s.getClipsFromDir(filepath.Join(videosDir, "fansly-clips"))
			allClips = append(allClips, videosClips...)
		}
	}

	// 3. Check custom directory if set
	if s.defaultOption == StoreInCustomDir && s.customOutputDir != "" {
		customClips := s.getClipsFromDir(s.customOutputDir)
		allClips = append(allClips, customClips...)
	}

	return allClips
}

// getClipsFromDir gets all video clips from a specific directory
func (s *ClipService) getClipsFromDir(dir string) []string {
	var clips []string

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return clips
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return clips
	}

	for _, file := range files {
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".mp4" || ext == ".webm" || ext == ".mov" {
				clips = append(clips, filepath.Join(dir, file.Name()))
			}
		}
	}

	return clips
}

// GetCurrentClipsDir returns the current directory where clips will be saved
func (s *ClipService) GetCurrentClipsDir(sourceVideoPath string) string {
	dir, _ := s.getOutputDirectory(sourceVideoPath)
	return dir
}

// OpenClipsFolder opens the current clips folder in the file explorer
func (s *ClipService) OpenClipsFolder(sourceVideoPath string) error {
	dir, err := s.getOutputDirectory(sourceVideoPath)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch os.Getenv("GOOS") {
	case "windows":
		cmd = exec.Command("explorer", dir)
	case "darwin":
		cmd = exec.Command("open", dir)
	default: // Linux and others
		cmd = exec.Command("xdg-open", dir)
	}

	return cmd.Start()
}

// Helper function to sanitize filenames
func sanitizeFilename(filename string) string {
	// Replace invalid characters with underscores
	invalidChars := []rune{'<', '>', ':', '"', '/', '\\', '|', '?', '*'}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, string(char), "_")
	}

	// Trim spaces from beginning and end
	filename = strings.TrimSpace(filename)

	// Ensure the filename isn't too long (max 255 chars)
	if len(filename) > 255 {
		filename = filename[:255]
	}

	// If filename is empty after sanitization, use a default name
	if filename == "" {
		filename = fmt.Sprintf("clip_%s", time.Now().Format("20060102_150405"))
	}

	return filepath.Clean(filename)
}

// Helper function to format time for ffmpeg
func formatFFmpegTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	milliseconds := int((seconds - float64(int(seconds))) * 1000)
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, milliseconds)
}
