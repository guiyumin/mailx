package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectorTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9370DB")).MarginBottom(1)
	selectorItemStyle   = lipgloss.NewStyle().PaddingLeft(4)
	selectorCursorStyle = lipgloss.NewStyle().PaddingLeft(4).Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#9370DB"))
	selectorHintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).MarginTop(1)
)

// SelectorItem represents an item in the selector
type SelectorItem struct {
	ID    string
	Label string
}

// Selector is a TUI component for selecting from a list
type Selector struct {
	title    string
	items    []SelectorItem
	cursor   int
	selected string
	cancelled bool
}

func NewSelector(title string, items []SelectorItem) Selector {
	return Selector{
		title: title,
		items: items,
	}
}

func (s Selector) Init() tea.Cmd {
	return nil
}

func (s Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.items)-1 {
				s.cursor++
			}
		case "enter":
			s.selected = s.items[s.cursor].ID
			return s, tea.Quit
		case "esc", "ctrl+c":
			s.cancelled = true
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s Selector) View() string {
	var b strings.Builder

	b.WriteString(selectorTitleStyle.Render(s.title))
	b.WriteString("\n\n")

	for i, item := range s.items {
		if i == s.cursor {
			b.WriteString(selectorCursorStyle.Render(fmt.Sprintf("> %s", item.Label)))
		} else {
			b.WriteString(selectorItemStyle.Render(fmt.Sprintf("  %s", item.Label)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(selectorHintStyle.Render("↑/k up • ↓/j down • enter select • esc cancel"))

	// Add padding to the whole view
	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (s Selector) Selected() string {
	return s.selected
}

func (s Selector) Cancelled() bool {
	return s.cancelled
}

// RunSelector runs the selector TUI and returns the selected ID and label
func RunSelector(title string, items []SelectorItem) (string, string, bool) {
	selector := NewSelector(title, items)
	p := tea.NewProgram(selector)

	m, err := p.Run()
	if err != nil {
		return "", "", true
	}

	result := m.(Selector)
	if result.Cancelled() {
		return "", "", true
	}

	// Find the label for the selected ID
	for _, item := range items {
		if item.ID == result.Selected() {
			return item.ID, item.Label, false
		}
	}

	return result.Selected(), "", false
}

// TextInput is a TUI component for text input
type TextInput struct {
	title     string
	input     textinput.Model
	cancelled bool
}

func NewTextInput(title, placeholder string) TextInput {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return TextInput{
		title: title,
		input: ti,
	}
}

func (t TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (t TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return t, tea.Quit
		case "esc", "ctrl+c":
			t.cancelled = true
			return t, tea.Quit
		}
	}

	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	return t, cmd
}

func (t TextInput) View() string {
	var b strings.Builder

	b.WriteString(selectorTitleStyle.Render(t.title))
	b.WriteString("\n\n")
	b.WriteString("    " + t.input.View())
	b.WriteString("\n\n")
	b.WriteString(selectorHintStyle.Render("enter confirm • esc cancel"))

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func (t TextInput) Value() string {
	return t.input.Value()
}

func (t TextInput) Cancelled() bool {
	return t.cancelled
}

// RunTextInput runs the text input TUI and returns the entered text
func RunTextInput(title, placeholder string) (string, bool) {
	ti := NewTextInput(title, placeholder)
	p := tea.NewProgram(ti)

	m, err := p.Run()
	if err != nil {
		return "", true
	}

	result := m.(TextInput)
	if result.Cancelled() || result.Value() == "" {
		return "", true
	}

	return result.Value(), false
}
