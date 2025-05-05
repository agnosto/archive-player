package services

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// FileDialogService handles file dialog operations
type FileDialogService struct {
	ctx context.Context
}

// NewFileDialogService creates a new file dialog service
func NewFileDialogService() *FileDialogService {
	return &FileDialogService{}
}

// SetContext sets the context for the service
func (s *FileDialogService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// OpenVideoFile opens a dialog to select a video file
func (s *FileDialogService) OpenVideoFile() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Select Video File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Video Files (*.mp4;*.webm;*.mkv;*.avi)",
				Pattern:     "*.mp4;*.webm;*.mkv;*.avi",
			},
		},
	}
	return runtime.OpenFileDialog(s.ctx, options)
}

// OpenChatFile opens a dialog to select a chat JSON file
func (s *FileDialogService) OpenChatFile() (string, error) {
	options := runtime.OpenDialogOptions{
		Title: "Select Chat JSON File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JSON Files (*.json)",
				Pattern:     "*.json",
			},
		},
	}
	return runtime.OpenFileDialog(s.ctx, options)
}
