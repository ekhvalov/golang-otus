package sender

import (
	"fmt"
	"io"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/notification"
)

func NewSender(w io.Writer) notification.Sender {
	return &s{w: w}
}

type s struct {
	w io.Writer
}

func (s s) Send(notification notification.Notification) error {
	_, err := s.w.Write([]byte(fmt.Sprintf(
		"%s %s %s\n",
		notification.EventTitle,
		notification.EventDate,
		notification.UserID,
	)))
	return err
}
