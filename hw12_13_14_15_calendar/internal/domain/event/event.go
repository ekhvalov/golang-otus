package event

import "time"

type Event struct {
	ID           string
	Title        string
	DateTime     time.Time
	Duration     time.Duration
	UserID       string
	Description  string
	NotifyBefore time.Duration
}

const MinDuration = time.Minute * 10
