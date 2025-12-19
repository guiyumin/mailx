# Bug: Search returns no results

## Problem
`./maily search -a danaoairuike@gmail.com -q "temu"` returned no results even though the user had many emails containing "temu" in their inbox.

## Root Cause
Gmail's IMAP implementation requires `UID SEARCH` instead of regular `SEARCH`. The go-imap v2 library's `client.Search()` method returned 0 results even when searching for ALL messages, while `client.UIDSearch()` works correctly.

## Debugging Steps

1. **Initial investigation**: The original code was using `X-GM-RAW` as a header field, which is incorrect. X-GM-RAW is a Gmail-specific IMAP extension for native search syntax.

2. **Tried TEXT search**: Changed to standard IMAP `Text` criteria - still 0 results.

3. **Tried BODY search**: Changed to `Body` criteria - still 0 results.

4. **Tried Subject header search**: Changed to `Header` with Subject key - still 0 results.

5. **Tested ALL search**: Added debug to search for ALL messages (empty criteria) - returned 0 UIDs even though mailbox had 8433 messages. This proved the `Search()` method itself was broken with Gmail.

6. **Switched to UIDSearch**: Changed from `client.Search()` to `client.UIDSearch()` - **this worked**.

## Fix

In `internal/gmail/imap.go`, changed:
```go
// Before (broken)
searchData, err := c.client.Search(searchCriteria, nil).Wait()

// After (working)
searchData, err := c.client.UIDSearch(searchCriteria, nil).Wait()
```

## Notes

- Gmail's IMAP has quirks compared to standard IMAP
- go-imap v2 doesn't support X-GM-RAW extension (would require raw command access)
- The search currently uses BODY criteria which searches message body content
- Search is performed on INBOX by default
