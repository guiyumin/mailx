# Mouse Scroll Jumping Issue

## Problem

When scrolling through the email list using the mouse wheel, the cursor would jump multiple emails at once (e.g., skip 2 emails and land on the 3rd) instead of moving one email at a time.

## Root Cause

Mouse wheel events in terminals often fire multiple times per physical scroll action. This is common with:
- Smooth scrolling trackpads
- High-sensitivity mice
- Terminal emulators that report fine-grained scroll events

A single scroll gesture could generate 3 or more scroll events, causing the cursor to move 3 positions instead of 1.

## Solution

Implemented count-based throttling that only processes every 3rd scroll event:

```go
case tea.MouseMsg:
    if a.state == stateReady && !a.confirmDelete {
        switch msg.Button {
        case tea.MouseButtonWheelUp:
            switch a.view {
            case listView:
                // Only process every 3rd scroll event
                a.scrollCount++
                if a.scrollCount >= 3 {
                    a.mailList.ScrollUp()
                    a.scrollCount = 0
                }
                return a, nil
            // ...
        }
    }
```

## Why Count-Based vs Time-Based Debouncing

Time-based debouncing (e.g., 30-50ms) was also tried but felt sluggish because:
- It introduces perceived lag between scroll gesture and response
- The "right" debounce time varies by hardware

Count-based throttling:
- Responds instantly to the Nth event (no perceived lag)
- Naturally adapts to the scroll hardware's event frequency
- Feels more responsive while still preventing jumps

## Files Changed

- `internal/ui/app.go`: Added `scrollCount int` field and count-based throttling in mouse event handler
