# AI Chat Feature Plan

## Overview

Add an AI chat feature to Maily that enables:
1. Email summarization and intelligent Q&A
2. Natural language calendar event creation
3. Context-aware assistance within the TUI

## Use Cases

### 1. Email Summarization
- Summarize current email thread
- Summarize all unread emails
- Extract action items from emails
- Answer questions about email content

### 2. Calendar Integration (Natural Language)

**CLI:**
```bash
maily c add "tomorrow 9am talk to Jerry"
maily c add "lunch with Sarah next Friday at noon"
maily c add "team standup every Monday 10am"
```

**TUI Event Add (hybrid approach):**
When pressing `a` to add event, user can choose:
```
┌─ Add Event ──────────────────────────────────────┐
│                                                  │
│  [1] Quick add (natural language)                │
│  [2] Manual input (form fields)                  │
│                                                  │
└──────────────────────────────────────────────────┘
```

Option 1 - Natural language:
```
┌─ Quick Add ──────────────────────────────────────┐
│                                                  │
│  > tomorrow 9am meeting with boss_               │
│                                                  │
│  Press Enter to parse, Esc to cancel             │
└──────────────────────────────────────────────────┘
```

Option 2 - Form fields (current UI):
```
┌─ Add Event ──────────────────────────────────────┐
│  Title: Meeting with boss                        │
│  Date:  2024-12-24                               │
│  Time:  09:00                                    │
│  Duration: 1h                                    │
└──────────────────────────────────────────────────┘
```

Both paths end at the same confirmation screen before saving.

### 3. Email-to-Calendar Extraction
- Detect dates/times/events mentioned in emails
- Offer to add detected events to calendar
- Examples:
  - "Let's meet Thursday at 3pm" → detected, prompt to add
  - "Deadline: Dec 31st" → detected, prompt to add reminder
  - Flight confirmations, restaurant reservations, etc.

**TUI Flow:**
```
┌─ Email from Jerry ─────────────────────────────┐
│ Hey, let's grab coffee tomorrow at 10am        │
│ at Blue Bottle on Market St.                   │
└────────────────────────────────────────────────┘

 [e] Add to calendar: "Coffee with Jerry - Tomorrow 10am"
```

Press `e` to extract and add to calendar with one keystroke.

### 4. TUI Chat Panel
- In-app chat for quick questions
- Context-aware (knows current email, calendar view)

### 5. Today View (Daily Dashboard)

Command: `maily today` or `maily t`

Split-panel view combining emails and events:

```
┌─ Today's Emails ─────────────────┬─ Events ──────────────┐
│ Meeting notes from Jerry         │ 9:00am                │
│ Q4 Budget Review                 │  Standup              │
│ Re: Project Timeline             │                       │
│ Invoice #1234                    │ 10:30am               │
│ Welcome to our newsletter        │  Meeting with boss    │
│                                  │                       │
│                                  │ 2:00pm                │
│                                  │  Client call          │
│                                  │                       │
│                                  │ 5:30pm                │
│                                  │  Gym                  │
└──────────────────────────────────┴───────────────────────┘
 [j/k] navigate  [enter] open  [a] add event  [e] edit  [d] delete
```

**Email Panel (Left):**
- Compact: title only (no date, no sender)
- Same navigation as full mail list (j/k, enter to open)
- Shows today's emails only (or unread?)

**Events Panel (Right):**
- Vertical timeline format
- Time on its own line, event title indented below
- Simple and scannable
- [a] add, [e] edit, [d] delete events

## Architecture

### CLI Commands

```bash
# Today view (daily dashboard)
maily today                          # or: maily t
                                     # Split view: emails + events

# Calendar shortcuts
maily c add "<natural language>"     # Add event via NLP
maily c list                         # List upcoming events

# Chat/AI commands
maily chat "<question>"              # One-shot question
maily chat                           # Enter interactive chat mode
```

### Components

```
internal/
├── ai/
│   ├── client.go          # AI provider abstraction (OpenAI, Anthropic, local)
│   ├── prompts.go         # System prompts for different tasks
│   ├── parser.go          # Parse NLP responses into structured data
│   └── context.go         # Build context from emails/calendar
├── calendar/
│   └── nlp.go             # Natural language date/time parsing
└── ui/
    └── components/
        ├── chatpanel.go   # TUI chat panel component
        ├── todayview.go   # Today dashboard (emails + events split)
        ├── compactmail.go # Compact email list (title only)
        └── eventlist.go   # Vertical event timeline
```

### AI Provider Strategy

**Target users:** CLI power users who already have AI tools installed.

**Reuse existing AI CLIs** (zero setup!):
- Claude Code: `claude -p "prompt" --output-format json`
- Codex: `codex exec "prompt" --json`
- Gemini: `gemini -p "prompt" --output-format json`
- Ollama: `ollama run llama3.2:3b "prompt"`

Auto-detect which CLI is available, just use it for everything:
- Date parsing ("tomorrow 9am" → structured event)
- Email summarization
- Event extraction from emails
- Any future AI features

### Natural Language Date Parsing

For `maily c add`, we need to parse natural language into structured event data:

```go
type ParsedEvent struct {
    Title     string
    StartTime time.Time
    EndTime   time.Time    // Optional, default 1 hour
    Recurrence string      // Optional: daily, weekly, etc.
}
```

**Approach:**
1. Send to AI: "Parse this into event details: tomorrow 9am talk to Jerry"
2. AI returns structured JSON
3. Create Google Calendar event via API

### Email Context Building

For summarization, build context from:
- Current email body + headers
- Thread history (if available)
- Sender information

## Implementation Phases

### Phase 0: Today View (No AI)
- [ ] Add `maily today` / `maily t` command
- [ ] Create compact email list component (title only)
- [ ] Create vertical event list component
- [ ] Build split-panel today view
- [ ] Add event CRUD (add/edit/delete) via keyboard
- [ ] Tab or arrow keys to switch between panels

### Phase 1: AI Integration (single phase for all AI features)
- [ ] Auto-detect available AI CLI (claude, codex, gemini, ollama)
- [ ] Implement `callAI(prompt) -> response` helper
- [ ] `maily c add "tomorrow 9am meeting"` - NLP event creation
- [ ] TUI quick-add with NLP (hybrid: quick add vs form)
- [ ] `e` key to extract events from current email
- [ ] `s` key to summarize current email
- [ ] `maily chat "question"` - one-shot Q&A

## Configuration

Add to `~/.config/maily/config.yml`:

```yaml
ai:
  provider: auto  # auto-detect, or: claude, codex, gemini, ollama
```

**Auto-detection order:** claude → codex → gemini → ollama

**Implementation:**
```go
func detectAI() string {
    for _, cli := range []string{"claude", "codex", "gemini", "ollama"} {
        if commandExists(cli) { return cli }
    }
    return "" // no AI available
}

func callAI(prompt string) (string, error) {
    switch detectAI() {
    case "claude":
        return exec.Command("claude", "-p", prompt, "--output-format", "json").Output()
    case "codex":
        return exec.Command("codex", "exec", prompt, "--json").Output()
    case "gemini":
        return exec.Command("gemini", "-p", prompt, "--output-format", "json").Output()
    case "ollama":
        return exec.Command("ollama", "run", "llama3.2:3b", prompt).Output()
    }
    return "", errors.New("no AI CLI found - install claude, codex, gemini, or ollama")
}
```

## Dependencies

- Existing Google Calendar API (from calendar feature)
- One of: Claude Code CLI, Codex CLI, Gemini CLI, or Ollama
