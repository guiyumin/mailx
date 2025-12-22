package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"maily/internal/calendar"
)

var (
	calendarID    string
	eventStart    string
	eventEnd      string
	eventLocation string
	eventNotes    string
	eventAllDay   bool
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Manage calendar events",
	Long:  `Access and manage calendar events from your macOS Calendar app.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default: show today's events
		showTodayEvents()
	},
}

var calendarListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all calendars",
	Run: func(cmd *cobra.Command, args []string) {
		listCalendars()
	},
}

var calendarEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Show upcoming events (next 7 days)",
	Run: func(cmd *cobra.Command, args []string) {
		showUpcomingEvents()
	},
}

var calendarAddCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new calendar event",
	Long: `Add a new event to your calendar.

Date/time format: "2024-01-15 10:00" or "2024-01-15" for all-day events`,
	Example: `  maily calendar add "Team Meeting" --start "2024-01-15 10:00" --end "2024-01-15 11:00"
  maily calendar add "Conference" --start "2024-01-20" --all-day
  maily calendar add "Lunch" --start "2024-01-15 12:00" --end "2024-01-15 13:00" --location "Cafe"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addEvent(args[0])
	},
}

var calendarDeleteCmd = &cobra.Command{
	Use:   "delete [event-id]",
	Short: "Delete a calendar event",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deleteEvent(args[0])
	},
}

func init() {
	calendarCmd.AddCommand(calendarListCmd)
	calendarCmd.AddCommand(calendarEventsCmd)
	calendarCmd.AddCommand(calendarAddCmd)
	calendarCmd.AddCommand(calendarDeleteCmd)

	calendarAddCmd.Flags().StringVarP(&eventStart, "start", "s", "", "Start time (required)")
	calendarAddCmd.Flags().StringVarP(&eventEnd, "end", "e", "", "End time (defaults to start + 1 hour)")
	calendarAddCmd.Flags().StringVarP(&calendarID, "calendar", "c", "", "Calendar ID (uses default if not specified)")
	calendarAddCmd.Flags().StringVarP(&eventLocation, "location", "l", "", "Event location")
	calendarAddCmd.Flags().StringVarP(&eventNotes, "notes", "n", "", "Event notes")
	calendarAddCmd.Flags().BoolVar(&eventAllDay, "all-day", false, "Create an all-day event")
	calendarAddCmd.MarkFlagRequired("start")
}

func getCalendarClient() (calendar.Client, error) {
	// Check status first to give better guidance
	status := calendar.GetAuthStatus()
	switch status {
	case calendar.AuthDenied:
		fmt.Println("Calendar access was denied.")
		fmt.Println()
		fmt.Println("To fix this:")
		fmt.Println("  1. Open System Settings > Privacy & Security > Calendars")
		fmt.Println("  2. Enable access for your terminal app")
		fmt.Println()
		return nil, calendar.ErrAccessDenied
	case calendar.AuthRestricted:
		fmt.Println("Calendar access is restricted by system policy.")
		return nil, calendar.ErrAccessDenied
	case calendar.AuthNotDetermined:
		fmt.Println("Requesting calendar access...")
	}

	client, err := calendar.NewClient()
	if err != nil {
		if err == calendar.ErrAccessDenied {
			fmt.Println()
			fmt.Println("Access was not granted. Please try again after enabling calendar access.")
		}
		return nil, err
	}
	return client, nil
}

func showTodayEvents() {
	client, err := getCalendarClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	events, err := client.ListEvents(startOfDay, endOfDay)
	if err != nil {
		fmt.Printf("Error loading events: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("No events today")
		return
	}

	fmt.Printf("Today's events (%s):\n\n", now.Format("Mon, Jan 2"))
	for _, event := range events {
		printEvent(event)
	}
}

func showUpcomingEvents() {
	client, err := getCalendarClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startOfDay.Add(7 * 24 * time.Hour)

	events, err := client.ListEvents(startOfDay, endDate)
	if err != nil {
		fmt.Printf("Error loading events: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("No upcoming events in the next 7 days")
		return
	}

	fmt.Println("Upcoming events (next 7 days):\n")

	var currentDate string
	for _, event := range events {
		eventDate := event.StartTime.Format("Mon, Jan 2")
		if eventDate != currentDate {
			if currentDate != "" {
				fmt.Println()
			}
			fmt.Printf("── %s ──\n", eventDate)
			currentDate = eventDate
		}
		printEvent(event)
	}
}

func printEvent(event calendar.Event) {
	var timeStr string
	if event.AllDay {
		timeStr = "All day"
	} else {
		timeStr = fmt.Sprintf("%s - %s",
			event.StartTime.Format("3:04 PM"),
			event.EndTime.Format("3:04 PM"))
	}

	fmt.Printf("  %s  %s\n", timeStr, event.Title)
	if event.Location != "" {
		fmt.Printf("            📍 %s\n", event.Location)
	}
	if event.Calendar != "" {
		fmt.Printf("            [%s]\n", event.Calendar)
	}
}

func listCalendars() {
	client, err := getCalendarClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	calendars, err := client.ListCalendars()
	if err != nil {
		fmt.Printf("Error loading calendars: %v\n", err)
		os.Exit(1)
	}

	if len(calendars) == 0 {
		fmt.Println("No calendars found")
		return
	}

	fmt.Println("Available calendars:\n")
	for _, cal := range calendars {
		fmt.Printf("  %s\n", cal.Title)
		fmt.Printf("    ID: %s\n", cal.ID)
	}
}

func parseDateTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)

	// Try full datetime format
	layouts := []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 3:04 PM",
		"2006-01-02 3:04PM",
		"2006-01-02",
		"01/02/2006 15:04",
		"01/02/2006",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date/time format: %s", s)
}

func addEvent(title string) {
	client, err := getCalendarClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	startTime, err := parseDateTime(eventStart)
	if err != nil {
		fmt.Printf("Error parsing start time: %v\n", err)
		os.Exit(1)
	}

	var endTime time.Time
	if eventEnd != "" {
		endTime, err = parseDateTime(eventEnd)
		if err != nil {
			fmt.Printf("Error parsing end time: %v\n", err)
			os.Exit(1)
		}
	} else if eventAllDay {
		// All-day events: end is same as start (or next day in some implementations)
		endTime = startTime.Add(24 * time.Hour)
	} else {
		// Default: 1 hour duration
		endTime = startTime.Add(1 * time.Hour)
	}

	event := calendar.Event{
		Title:     title,
		StartTime: startTime,
		EndTime:   endTime,
		Calendar:  calendarID,
		Location:  eventLocation,
		Notes:     eventNotes,
		AllDay:    eventAllDay,
	}

	eventID, err := client.CreateEvent(event)
	if err != nil {
		fmt.Printf("Error creating event: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Event created: %s\n", title)
	fmt.Printf("  Time: %s - %s\n", startTime.Format("Mon, Jan 2 3:04 PM"), endTime.Format("3:04 PM"))
	if eventLocation != "" {
		fmt.Printf("  Location: %s\n", eventLocation)
	}
	fmt.Printf("  ID: %s\n", eventID)
}

func deleteEvent(eventID string) {
	client, err := getCalendarClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = client.DeleteEvent(eventID)
	if err != nil {
		if err == calendar.ErrNotFound {
			fmt.Printf("Event not found: %s\n", eventID)
		} else {
			fmt.Printf("Error deleting event: %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("Event deleted: %s\n", eventID)
}
