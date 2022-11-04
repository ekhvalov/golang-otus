package command

import (
	"fmt"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
)

type ErrValidate struct {
	message string
}

func (e ErrValidate) Error() string {
	return e.message
}

func NewErrValidate(message string) ErrValidate {
	return ErrValidate{message: message}
}

func validateTitle(title string) error {
	if title == "" {
		return NewErrValidate("title is empty")
	}
	return nil
}

func validateDateTime(dateTime time.Time) error {
	if dateTime.Before(time.Now()) {
		return NewErrValidate("requested date is in the past")
	}
	return nil
}

func validateDuration(duration time.Duration) error {
	if duration < event.MinDuration {
		return NewErrValidate(fmt.Sprintf("minimal duration is %s", event.MinDuration))
	}
	return nil
}

func validateUserID(userID string) error {
	if userID == "" {
		return NewErrValidate("user ID is empty")
	}
	return nil
}

func validateID(id string) error {
	if id == "" {
		return NewErrValidate("ID is empty")
	}
	return nil
}
