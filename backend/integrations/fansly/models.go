package fansly

import ()

// Config represents the configuration for Fansly integration
type Config struct {
	ConfigPath string `json:"configPath"`
	DbPath     string `json:"dbPath"`
}

// Stream represents a stream from the Fansly database
type Stream struct {
	Model        string  `json:"model"`
	Hash         string  `json:"hash"`
	Path         string  `json:"path"`
	FileType     string  `json:"file_type"`
	ContactSheet string  `json:"contactSheet,omitempty"`
	Duration     float64 `json:"duration,omitempty"`
}

// StreamsResult represents the result of getting Fansly streams
type StreamsResult struct {
	Streams   []Stream `json:"streams"`
	ChatFiles []string `json:"chatFiles"`
	Error     string   `json:"error,omitempty"`
}

// StreamResult represents the result of loading a Fansly stream
type StreamResult struct {
	VideoPath    string `json:"videoPath"`
	ChatPath     string `json:"chatPath"`
	ContactSheet string `json:"contactSheet"`
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
}
