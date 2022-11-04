package notification

type Sender interface {
	// Send a Notification to a user
	Send(notification Notification) error
}
