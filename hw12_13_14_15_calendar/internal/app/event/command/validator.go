package command

import (
	"fmt"
	"time"

	"github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/internal/domain/event"
)

func validateTitle(title string) error {
	if title == "" {
		return fmt.Errorf("title is empty")
	}
	return nil
}

func validateDateTime(dateTime time.Time) error {
	if dateTime.Before(time.Now()) {
		return fmt.Errorf("requested date is in the past")
	}
	return nil
}

func validateDuration(duration time.Duration) error {
	if duration < event.MinDuration {
		return fmt.Errorf("minimal duration is: %s", event.MinDuration)
	}
	return nil
}

func validateUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is empty")
	}
	return nil
}

func validateID(id string) error {
	if id == "" {
		return fmt.Errorf("ID is empty")
	}
	return nil
}
