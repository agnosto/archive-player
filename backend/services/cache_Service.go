package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// VideoMetadata represents cached information about a video file
type VideoMetadata struct {
	Path         string    `json:"path"`
	Hash         string    `json:"hash"`
	Duration     float64   `json:"duration"`
	LastModified time.Time `json:"lastModified"`
	FileSize     int64     `json:"fileSize"`
}

// VideoCache represents the cache of video metadata
type VideoCache struct {
	Videos      map[string]VideoMetadata `json:"videos"` // key is video path
	LastUpdated time.Time                `json:"lastUpdated"`
}

// CacheService handles caching of video metadata
type CacheService struct {
	cacheDir string
}

// NewCacheService creates a new cache service
func NewCacheService(cacheDir string) *CacheService {
	return &CacheService{
		cacheDir: cacheDir,
	}
}

// GetCachePath returns the path to the cache file
func (s *CacheService) GetCachePath(cacheType string) string {
	return filepath.Join(s.cacheDir, cacheType+"_cache.json")
}

// LoadVideoCache loads the video cache from disk
func (s *CacheService) LoadVideoCache(cacheType string) (VideoCache, error) {
	cache := VideoCache{
		Videos:      make(map[string]VideoMetadata),
		LastUpdated: time.Time{},
	}

	cachePath := s.GetCachePath(cacheType)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Cache doesn't exist yet, return empty cache
			return cache, nil
		}
		return cache, err
	}

	err = json.Unmarshal(data, &cache)
	if err != nil {
		// If cache is corrupted, return empty cache
		return VideoCache{
			Videos:      make(map[string]VideoMetadata),
			LastUpdated: time.Time{},
		}, nil
	}

	return cache, nil
}

// SaveVideoCache saves the video cache to disk
func (s *CacheService) SaveVideoCache(cacheType string, cache VideoCache) error {
	cache.LastUpdated = time.Now()

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	cachePath := s.GetCachePath(cacheType)

	// Create directory if it doesn't exist
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}
