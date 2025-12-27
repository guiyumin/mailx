package utils

import (
	"strings"
	"time"
)

// Helper functions
func TruncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func ExtractNameFromEmail(from string) string {
	if idx := strings.Index(from, "<"); idx > 0 {
		return strings.TrimSpace(from[:idx])
	}
	return from
}

func FormatEmailDate(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04")
	}
	if t.Year() == now.Year() {
		return t.Format("Jan 02")
	}
	return t.Format("02/01/06")
}
