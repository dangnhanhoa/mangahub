package models

import "time"

// ReadingStatus values
const (
	StatusReading    = "reading"
	StatusCompleted  = "completed"
	StatusPlanToRead = "plan-to-read"
	StatusOnHold     = "on-hold"
	StatusDropped    = "dropped"
)

type UserProgress struct {
	UserID         string    `json:"user_id" db:"user_id"`
	MangaID        string    `json:"manga_id" db:"manga_id"`
	CurrentChapter int       `json:"current_chapter" db:"current_chapter"`
	Status         string    `json:"status" db:"status"`
	Rating         int       `json:"rating,omitempty" db:"rating"` // 0 = unrated, 1–10
	Notes          string    `json:"notes,omitempty" db:"notes"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ProgressUpdate is broadcast over TCP when a user updates their progress.
type ProgressUpdate struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	MangaID   string `json:"manga_id"`
	MangaTitle string `json:"manga_title,omitempty"`
	Chapter   int    `json:"chapter"`
	Timestamp int64  `json:"timestamp"`
}

// Notification is broadcast over UDP.
type Notification struct {
	Type      string `json:"type"`    // "chapter_update" | "system" | "test"
	MangaID   string `json:"manga_id,omitempty"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// ChatMessage is sent over WebSocket.
type ChatMessage struct {
	Type      string `json:"type"` // "message" | "join" | "leave" | "system"
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func ValidReadingStatus(s string) bool {
	switch s {
	case StatusReading, StatusCompleted, StatusPlanToRead, StatusOnHold, StatusDropped:
		return true
	}
	return false
}
