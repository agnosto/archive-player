package integrations

import (
	"FanslyArchivePlayer/backend/integrations/fansly"
	"FanslyArchivePlayer/backend/services"
)

// Manager handles all integrations
type Manager struct {
	FanslyService *fansly.Service
}

// NewManager creates a new integrations manager
func NewManager(appDataDir string, cacheService *services.CacheService) *Manager {
	return &Manager{
		FanslyService: fansly.NewService(appDataDir, cacheService),
	}
}
