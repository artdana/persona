package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmModel struct {
	question string
	message  string
	selected bool
	confirmed bool
	cancelled bool
}

func NewConfirmModel(question, message string) ConfirmModel {
	return ConfirmModel{
		question:  question,
		message:   message,
		selected:  false,
		confirmed: false,
		cancelled: false,
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "left", "h":
			m.selected = false
		case "right", "l":
			m.selected = true
		case "enter":
			m.confirmed = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	var b strings.Builder

	// Question
	b.WriteString(TitleStyle.Render(m.question))
	b.WriteString("\n\n")

	// Message
	if m.message != "" {
		b.WriteString(ValueStyle.Render(m.message))
		b.WriteString("\n\n")
	}

	// Buttons
	yesStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(MutedColor)

	noStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(MutedColor)

	if m.selected {
		yesStyle = yesStyle.
			BorderForeground(FocusedColor).
			Foreground(FocusedColor)
	} else {
		noStyle = noStyle.
			BorderForeground(FocusedColor).
			Foreground(FocusedColor)
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		noStyle.Render("No"),
		"  ",
		yesStyle.Render("Yes"),
	)

	b.WriteString(buttons)
	b.WriteString("\n\n")

	// Help
	help := "←/→ or h/l: Switch  Enter: Confirm  Esc: Cancel"
	b.WriteString(HelpStyle.Render(help))

	return BorderStyle.Render(b.String())
}

func (m ConfirmModel) Confirmed() bool {
	return m.confirmed && m.selected
}

func (m ConfirmModel) Cancelled() bool {
	return m.cancelled
}

