package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"FanslyArchivePlayer/backend/models"
)

// VideoService handles video and chat operations
type VideoService struct {
	CurrentVideoPath string
	ChatMessages     []models.ChatMessage
}

// NewVideoService creates a new video service
func NewVideoService() *VideoService {
	return &VideoService{
		ChatMessages: []models.ChatMessage{},
	}
}

func extractTipAmount(rawData string) (int, bool) {
	if rawData == "" {
		return 0, false
	}

	// Try to parse the raw data
	var rawDataObj map[string]interface{}
	err := json.Unmarshal([]byte(rawData), &rawDataObj)
	if err != nil {
		return 0, false
	}

	// Check if event exists
	eventStr, ok := rawDataObj["event"].(string)
	if !ok {
		return 0, false
	}

	// Parse the event string
	var eventObj map[string]interface{}
	err = json.Unmarshal([]byte(eventStr), &eventObj)
	if err != nil {
		return 0, false
	}

	// Check if chatRoomMessage exists
	chatRoomMsg, ok := eventObj["chatRoomMessage"].(map[string]interface{})
	if !ok {
		return 0, false
	}

	// Check if attachments exist
	attachments, ok := chatRoomMsg["attachments"].([]interface{})
	if !ok || len(attachments) == 0 {
		return 0, false
	}

	// Look for metadata with amount in attachments
	for _, attachment := range attachments {
		attachmentMap, ok := attachment.(map[string]interface{})
		if !ok {
			continue
		}

		metadataStr, ok := attachmentMap["metadata"].(string)
		if !ok {
			continue
		}

		var metadata map[string]interface{}
		err = json.Unmarshal([]byte(metadataStr), &metadata)
		if err != nil {
			continue
		}

		// Check if amount exists
		amount, ok := metadata["amount"].(float64)
		if ok {
			return int(amount), true
		}
	}

	return 0, false
}

// LoadVideo loads a video file
func (s *VideoService) LoadVideo(path string) error {
	// Check if file exists
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("video file not found: %v", err)
	}

	s.CurrentVideoPath = path
	return nil
}

// LoadChatFile loads a chat JSON file
func (s *VideoService) LoadChatFile(path string) error {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read chat file: %v", err)
	}

	// Try to parse as an array of messages first
	var messages []models.ChatMessage
	err = json.Unmarshal(data, &messages)

	// If that fails, try to parse as a ChatData object with a messages field
	if err != nil {
		var chatData models.ChatData
		err = json.Unmarshal(data, &chatData)
		if err != nil {
			// Try to parse as a single message (and wrap it in an array)
			var singleMessage models.ChatMessage
			err = json.Unmarshal(data, &singleMessage)
			if err != nil {
				return fmt.Errorf("failed to parse chat JSON: %v", err)
			}
			messages = []models.ChatMessage{singleMessage}
		} else {
			messages = chatData.Messages
		}
	}

	// If we still have no messages, return an error
	if len(messages) == 0 {
		return fmt.Errorf("no chat messages found in file")
	}

	// Process each message to extract tip amount if present
	for i := range messages {
		if tipAmount, found := extractTipAmount(messages[i].RawData); found {
			messages[i].TipAmount = tipAmount
		}
	}

	// Sort messages by timestamp
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].TimeInSeconds < messages[j].TimeInSeconds
	})

	s.ChatMessages = messages
	return nil
}

// GetMessagesAtTime returns messages within a time window
func (s *VideoService) GetMessagesAtTime(currentTime float64, windowSize float64) []models.ChatMessage {
	if len(s.ChatMessages) == 0 {
		return []models.ChatMessage{}
	}

	// Find messages within the time window
	var messagesInWindow []models.ChatMessage
	for _, msg := range s.ChatMessages {
		if msg.TimeInSeconds >= currentTime-windowSize && msg.TimeInSeconds <= currentTime {
			messagesInWindow = append(messagesInWindow, msg)
		}
	}

	return messagesInWindow
}

// GetVideoFileInfo returns information about the current video
func (s *VideoService) GetVideoFileInfo() map[string]string {
	info := make(map[string]string)

	if s.CurrentVideoPath == "" {
		return info
	}

	info["path"] = s.CurrentVideoPath
	info["filename"] = filepath.Base(s.CurrentVideoPath)

	// Try to find associated files
	basePath := strings.TrimSuffix(s.CurrentVideoPath, filepath.Ext(s.CurrentVideoPath))

	// Check for chat file
	chatPath := basePath + "_chat.json"
	if _, err := os.Stat(chatPath); err == nil {
		info["chatFile"] = chatPath
	}

	// Check for thumbnail
	thumbnailPath := basePath + "_contact_sheet.jpg"
	if _, err := os.Stat(thumbnailPath); err == nil {
		info["thumbnailPath"] = thumbnailPath
	}

	return info
}
