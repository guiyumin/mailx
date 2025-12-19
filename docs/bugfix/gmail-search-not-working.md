# Bug: Gmail Search returns no results

## Problem
`./maily search -a danaoairuike@gmail.com -q "temu"` returned no results even though the user had many emails containing "temu" in their inbox.

## Root Cause
Two issues:

1. **go-imap v2's `Search()` doesn't work with Gmail** - Even searching for ALL messages returned 0 UIDs. Only `UIDSearch()` works.

2. **go-imap v2 doesn't support X-GM-RAW** - The library doesn't expose raw command access, so we can't use Gmail's powerful search syntax (`from:`, `has:attachment`, `category:`, etc.).

## Debugging Steps

1. **Initial investigation**: Original code used `X-GM-RAW` as a header field - incorrect usage.

2. **Tried TEXT search**: Changed to standard IMAP `Text` criteria - 0 results.

3. **Tried BODY search**: Changed to `Body` criteria - 0 results.

4. **Tried Subject header search**: Changed to `Header` with Subject key - 0 results.

5. **Tested ALL search**: Empty criteria returned 0 UIDs even with 8433 messages in mailbox. Proved `Search()` was broken with Gmail.

6. **Switched to UIDSearch**: `client.UIDSearch()` worked, but still limited to standard IMAP criteria.

7. **Final solution**: Bypass go-imap entirely for search using raw IMAP commands.

## Final Solution

Created `internal/gmail/search.go` - a raw IMAP implementation that:

1. Opens direct TLS connection to Gmail's IMAP server
2. Sends raw IMAP commands: LOGIN, SELECT, UID SEARCH X-GM-RAW
3. Parses the response to extract UIDs
4. Returns UIDs to be fetched via go-imap's `Fetch()`

```go
// Raw IMAP command sent to Gmail
searchCmd := fmt.Sprintf("a3 UID SEARCH X-GM-RAW %s\r\n", quoteString(query))
```

This enables Gmail's full search syntax:
- `temu` - simple text search
- `from:temu` - search by sender
- `category:promotions older_than:30d` - promotional emails older than 30 days
- `has:attachment is:unread` - unread emails with attachments

## Architecture

```
Search Flow:
┌──────────────────────────────────────────────────────┐
│ search.go (raw IMAP)                                 │
│   - TLS connection to imap.gmail.com:993             │
│   - LOGIN, SELECT, UID SEARCH X-GM-RAW "query"       │
│   - Returns: []UID                                   │
└──────────────────────────────────────────────────────┘
                          │
                          ▼
┌──────────────────────────────────────────────────────┐
│ imap.go (go-imap v2)                                 │
│   - client.Fetch(uidSet, fetchOptions)               │
│   - Returns: []Email with full details               │
└──────────────────────────────────────────────────────┘
```

## Why Not Extend go-imap v2?

- `SearchCriteria` struct has no field for raw/custom search keys
- `writeSearchKey()` function is internal, can't be extended
- `Client.beginCommand()` is unexported, no raw command access
- Would require forking the entire library

## Files Changed

- `internal/gmail/search.go` - New file: raw IMAP Gmail search with X-GM-RAW
- `internal/gmail/imap.go` - Updated `SearchMessages()` to use `GmailSearch()`
- `internal/cli/search.go` - Updated help text with Gmail search syntax examples
