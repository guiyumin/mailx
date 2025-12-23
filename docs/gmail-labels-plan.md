# Gmail Labels/Folders Support

## Overview

Maily supports viewing emails from any Gmail label or folder, not just Inbox.

**Important**: Maily is read-only for labels. We don't create, rename, or delete labels - that's Gmail's job. We just let users switch between existing labels to view their emails.

## How It Works

1. **Mail list** displays emails from the currently selected label
2. **Press `g`** to open the label picker
3. **Select a label** to view its emails
4. **Header badge** shows current label when not in Inbox

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `g` | Open label picker |
| `↑/↓` or `j/k` | Navigate labels |
| `Enter` | Select label |
| `Esc` | Cancel |

## Label Display

Gmail IMAP labels are shown with friendly names:

| IMAP Name | Display Name |
|-----------|--------------|
| `INBOX` | Inbox |
| `[Gmail]/Sent Mail` | Sent |
| `[Gmail]/Drafts` | Drafts |
| `[Gmail]/Spam` | Spam |
| `[Gmail]/Trash` | Trash |
| `[Gmail]/All Mail` | All Mail |
| `[Gmail]/Starred` | Starred |
| `[Gmail]/Important` | Important |
| Custom labels | As-is |

## Implementation

### Files

| File | Purpose |
|------|---------|
| `internal/ui/components/labelpicker.go` | LabelPicker component |
| `internal/ui/app.go` | Label state, picker toggle, label switching |
| `internal/ui/commands.go` | `loadLabels()`, `loadEmails()` uses current label |
| `internal/ui/components/views.go` | Header shows current label badge |

### Key Components

**LabelPicker** (`components.LabelPicker`)
- Full-screen modal for selecting labels
- Shows system labels first (sorted), then custom labels
- Indicates currently active label with `●`

**App State**
- `labelPicker` - the picker component
- `currentLabel` - currently viewing label (default "INBOX")
- `showLabelPicker` - whether picker is visible
- `emailCache` - keyed by `"accountIdx:label"` for per-label caching

### Flow

```
1. App starts -> loads labels via ListMailboxes()
2. User presses 'g' -> showLabelPicker = true
3. User navigates and selects -> currentLabel updated
4. loadEmails() fetches from currentLabel
5. Search also uses currentLabel
```
