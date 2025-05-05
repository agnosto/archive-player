package models

import "time"

// Author represents the author of a chat message
type Author struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Images   []string  `json:"images,omitempty"`
	Badges   []string  `json:"badges,omitempty"`
	TierInfo *TierInfo `json:"tier_info,omitempty"`
}

// TierInfo represents subscription tier information
type TierInfo struct {
	TierID    string `json:"tier_id"`
	TierColor string `json:"tier_color"`
	TierName  string `json:"tier_name"`
}

// ChatMessage represents a single chat message
type ChatMessage struct {
	MessageID     string    `json:"message_id"`
	Message       string    `json:"message"`
	MessageType   string    `json:"message_type"`
	Timestamp     int64     `json:"timestamp"`
	TimeInSeconds float64   `json:"time_in_seconds"`
	TimeText      string    `json:"time_text"`
	Author        Author    `json:"author"`
	RawData       string    `json:"raw_data,omitempty"`
	ReceivedAt    time.Time `json:"received_at,omitempty"`
	TipAmount     int       `json:"tip_amount,omitempty"`
}

// ChatData represents a collection of chat messages
type ChatData struct {
	Messages []ChatMessage `json:"messages"`
}
