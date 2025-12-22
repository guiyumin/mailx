#ifndef EVENTKIT_DARWIN_H
#define EVENTKIT_DARWIN_H

// Result codes
#define EK_SUCCESS 0
#define EK_ERROR_ACCESS_DENIED 1
#define EK_ERROR_NOT_FOUND 2
#define EK_ERROR_FAILED 3

// Authorization status codes
#define EK_AUTH_NOT_DETERMINED 0
#define EK_AUTH_RESTRICTED 1
#define EK_AUTH_DENIED 2
#define EK_AUTH_AUTHORIZED 3

// Check current authorization status without prompting
// Returns one of EK_AUTH_* codes
int GetAuthorizationStatus(void);

// Request calendar access from the user (triggers dialog if not determined)
// Returns EK_SUCCESS if granted, EK_ERROR_ACCESS_DENIED if denied
int RequestCalendarAccess(void);

// List all calendars
// Returns JSON array of calendars: [{"id":"...", "title":"...", "color":"..."}]
// Caller must free the returned string
char* ListCalendars(void);

// List events between start and end dates (Unix timestamps)
// Returns JSON array of events
// Caller must free the returned string
char* ListEvents(long long startTimestamp, long long endTimestamp);

// Create a new event
// Returns the event ID on success, NULL on failure
// Caller must free the returned string
char* CreateEvent(const char* title, long long startTimestamp, long long endTimestamp,
                  const char* calendarID, const char* location, const char* notes, int allDay);

// Delete an event by ID
// Returns EK_SUCCESS on success, error code on failure
int DeleteEvent(const char* eventID);

// Free a string returned by the EventKit functions
void FreeString(char* str);

#endif // EVENTKIT_DARWIN_H
