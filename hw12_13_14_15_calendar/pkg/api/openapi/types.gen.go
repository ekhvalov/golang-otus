// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package openapi

// Defines values for EventsPeriod.
const (
	Day   EventsPeriod = "day"
	Month EventsPeriod = "month"
	Week  EventsPeriod = "week"
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// Event defines model for Event.
type Event struct {
	// Event start date (unix timestamp)
	Date int64 `json:"date"`

	// Event description
	Description *string `json:"description,omitempty"`

	// Duration of an event (in minutes)
	Duration int    `json:"duration"`
	Id       string `json:"id"`

	// Amount of minutes to notify user before the event
	NotifyBefore *int `json:"notifyBefore,omitempty"`

	// Event title
	Title string `json:"title"`
}

// EventsPeriod defines model for EventsPeriod.
type EventsPeriod string

// NewEvent defines model for NewEvent.
type NewEvent struct {
	// Event start date (unix timestamp)
	Date int64 `json:"date"`

	// Event description
	Description *string `json:"description,omitempty"`

	// Duration of an event (in minutes)
	Duration int `json:"duration"`

	// Amount of minutes to notify user before the event
	NotifyBefore *int `json:"notifyBefore,omitempty"`

	// Event title
	Title string `json:"title"`
}

// GetEventsParams defines parameters for GetEvents.
type GetEventsParams struct {
	Period EventsPeriod `form:"period" json:"period"`
	Date   int64        `form:"date" json:"date"`
}

// PostEventsJSONBody defines parameters for PostEvents.
type PostEventsJSONBody = NewEvent

// PutEventsJSONBody defines parameters for PutEvents.
type PutEventsJSONBody = Event

// PostEventsJSONRequestBody defines body for PostEvents for application/json ContentType.
type PostEventsJSONRequestBody = PostEventsJSONBody

// PutEventsJSONRequestBody defines body for PutEvents for application/json ContentType.
type PutEventsJSONRequestBody = PutEventsJSONBody