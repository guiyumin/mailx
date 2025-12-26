package calendar

import (
	"errors"
	"time"
)

var (
	ErrAccessDenied = errors.New("calendar access denied")
	ErrNotFound     = errors.New("event not found")
	ErrFailed       = errors.New("operation failed")
	ErrNotSupported = errors.New("calendar integration only supported on macOS")
)

// Event represents a calendar event
type Event struct {
	ID                 string
	Title              string
	StartTime          time.Time
	EndTime            time.Time
	Location           string
	Notes              string
	Calendar           string
	AllDay             bool
	AlarmMinutesBefore int // Minutes before event to trigger alarm (0 = no alarm)
}

// Calendar represents a calendar source
type Calendar struct {
	ID    string
	Title string
	Color string
}

// Client provides access to calendar operations
type Client interface {
	// ListCalendars returns all available calendars
	ListCalendars() ([]Calendar, error)

	// ListEvents returns events within the given time range
	ListEvents(start, end time.Time) ([]Event, error)

	// CreateEvent creates a new event and returns its ID
	CreateEvent(event Event) (string, error)

	// DeleteEvent removes an event by its ID
	DeleteEvent(id string) error
}
