package notification

import "time"

type Notification struct {
	EventID    string
	EventTitle string
	EventDate  time.Time
	UserID     string
}
