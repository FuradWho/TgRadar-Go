package model

import "time"

// MessageData holds raw message info
type MessageData struct {
	GroupID   int64
	SenderID  int64
	Text      string
	Timestamp time.Time
}

type GroupStats struct {
	MsgCount  int
	UserCount int
	TopWords  []string
}

// AnalysisResult holds the AI analysis output
type AnalysisResult struct {
	Summary   string
	Sentiment string // Positive/Negative/Neutral
	HotTopics []string
}
