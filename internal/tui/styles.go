package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	PrimaryColor   = lipgloss.Color("63")   // Purple
	SuccessColor   = lipgloss.Color("10")   // Green
	WarningColor   = lipgloss.Color("11")   // Yellow
	ErrorColor     = lipgloss.Color("9")    // Red
	FocusedColor   = lipgloss.Color("12")   // Blue
	MutedColor     = lipgloss.Color("240")  // Gray
	BorderColor    = lipgloss.Color("62")   // Purple border

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			MarginBottom(1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginRight(2)

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// Interactive styles
	ActiveMarkerStyle = "✅"
	CursorStyle       = "➜"

	FocusedStyle = lipgloss.NewStyle().
			Foreground(FocusedColor)

	ActiveStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	// Border styles
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginTop(1)
)

