package components

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors
var (
	Primary   = lipgloss.Color("#7C3AED")
	Secondary = lipgloss.Color("#A78BFA")
	Success   = lipgloss.Color("#10B981")
	Warning   = lipgloss.Color("#F59E0B")
	Danger    = lipgloss.Color("#EF4444")
	Muted     = lipgloss.Color("#6B7280")
	Text      = lipgloss.Color("#F9FAFB")
	TextDim   = lipgloss.Color("#9CA3AF")
	Bg        = lipgloss.Color("#1F2937")
	BgDark    = lipgloss.Color("#111827")
)

// Base styles
var (
	BaseStyle = lipgloss.NewStyle().
			Background(BgDark).
			Foreground(Text)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Padding(0, 1).
			MarginBottom(1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text).
			Background(Primary).
			Padding(0, 2)

	// List styles
	ListStyle = lipgloss.NewStyle().
			Padding(1, 2)

	SelectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(Text).
				Background(Primary).
				Padding(0, 1)

	NormalItemStyle = lipgloss.NewStyle().
			Foreground(Text).
			Padding(0, 1)

	UnreadItemStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text).
			Padding(0, 1)

	// Email preview styles
	PreviewStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Muted).
			Padding(1, 2).
			MarginLeft(2)

	FromStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Secondary)

	SubjectStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text)

	DateStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Italic(true)

	SnippetStyle = lipgloss.NewStyle().
			Foreground(TextDim)

	// Status bar styles
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(TextDim).
			Background(Bg).
			Padding(0, 1)

	StatusKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			Background(Bg).
			Padding(0, 1)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
			Foreground(Muted).
			Padding(1, 2)

	HelpKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Secondary)

	HelpDescStyle = lipgloss.NewStyle().
			Foreground(Muted)

	// Loading/spinner styles
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(Primary)

	// Error styles
	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Danger).
			Padding(1, 2)

	// Success styles
	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Success).
			Padding(1, 2)

	// Tab styles
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text).
			Background(Primary).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(Muted).
				Background(Bg).
				Padding(0, 2)

	// Badge styles
	UnreadBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(Text).
			Background(Primary).
			Padding(0, 1)

	LabelBadge = lipgloss.NewStyle().
			Foreground(Text).
			Background(Muted).
			Padding(0, 1)

	// Dialog styles
	DialogStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 3).
			Align(lipgloss.Center)

	DialogTitleStyle = lipgloss.NewStyle().
				Bold(true)

	DialogHintStyle = lipgloss.NewStyle().
			Foreground(TextDim)
)
