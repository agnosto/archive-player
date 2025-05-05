package services

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// StreamService handles video streaming
type StreamService struct{}

// NewStreamService creates a new stream service
func NewStreamService() *StreamService {
	return &StreamService{}
}

// ServeVideoFile serves a video file with support for range requests
func (s *StreamService) ServeVideoFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Could not open video file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not get file info", http.StatusInternalServerError)
		return
	}

	fileSize := fileInfo.Size()
	contentType := getContentTypeFromExtension(filepath.Ext(filePath))

	// Set content type
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")

	// Check if Range header exists
	rangeHeader := r.Header.Get("Range")
	if rangeHeader != "" {
		// Parse the Range header
		start, end, err := parseRange(rangeHeader, fileSize)
		if err != nil {
			http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
			return
		}

		// Set the appropriate headers for a partial response
		w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
		w.Header().Set("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(fileSize, 10))
		w.WriteHeader(http.StatusPartialContent)

		// Seek to the start position
		_, err = file.Seek(start, io.SeekStart)
		if err != nil {
			http.Error(w, "Failed to seek file", http.StatusInternalServerError)
			return
		}

		// Create a limited reader to read only the requested range
		limitedReader := io.LimitReader(file, end-start+1)
		io.Copy(w, limitedReader)
	} else {
		// No range requested, serve the entire file
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		io.Copy(w, file)
	}
}

// Helper function to parse the Range header
func parseRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	var start, end int64
	var err error

	// Parse the range header
	_, rangeValue := splitHeaderValue(rangeHeader)
	if rangeValue == "" {
		return 0, fileSize - 1, nil
	}

	// Remove "bytes=" prefix
	if len(rangeValue) >= 6 && rangeValue[:6] == "bytes=" {
		rangeValue = rangeValue[6:]
	}

	// Split the range value
	rangeParts := splitRange(rangeValue)
	if len(rangeParts) != 2 {
		return 0, 0, err
	}

	// Parse start value
	if rangeParts[0] != "" {
		start, err = strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	// Parse end value
	if rangeParts[1] != "" {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	} else {
		end = fileSize - 1
	}

	// Validate range
	if start >= fileSize {
		return 0, 0, err
	}
	if end >= fileSize {
		end = fileSize - 1
	}

	return start, end, nil
}

// Helper function to split header value
func splitHeaderValue(header string) (string, string) {
	for i := 0; i < len(header); i++ {
		if header[i] == ':' {
			return header[:i], header[i+1:]
		}
	}
	return header, ""
}

// Helper function to split range value
func splitRange(rangeValue string) []string {
	result := make([]string, 2)
	for i := 0; i < len(rangeValue); i++ {
		if rangeValue[i] == '-' {
			result[0] = rangeValue[:i]
			result[1] = rangeValue[i+1:]
			return result
		}
	}
	result[0] = rangeValue
	return result
}

// Helper function to get content type from file extension
func getContentTypeFromExtension(ext string) string {
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	default:
		return "video/mp4"
	}
}
