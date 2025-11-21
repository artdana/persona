package tui

import (
	"persona/internal/persona"

	tea "github.com/charmbracelet/bubbletea"
)

func StartProfileSelector(profiles []persona.Profile, activeName string) *persona.Profile {
	p := tea.NewProgram(NewProfileModel(profiles, activeName), tea.WithAltScreen())
	m, _ := p.Run()
	finalModel := m.(Model)
	return finalModel.SelectedProfile()
}

// StartForm starts a form TUI and returns the values or nil if cancelled
func StartForm(title string, fields []FormField, onSubmit func(map[string]string) bool) map[string]string {
	model := NewFormModel(title, fields, onSubmit)
	p := tea.NewProgram(model, tea.WithAltScreen())
	m, _ := p.Run()
	finalModel := m.(FormModel)
	if finalModel.Cancelled() {
		return nil
	}
	return finalModel.Values()
}

// StartConfirm starts a confirmation dialog TUI
func StartConfirm(question, message string) bool {
	model := NewConfirmModel(question, message)
	p := tea.NewProgram(model, tea.WithAltScreen())
	m, _ := p.Run()
	finalModel := m.(ConfirmModel)
	return finalModel.Confirmed()
}

// StartList starts a list view TUI
func StartList(profiles []persona.Profile, activeName, title string) {
	model := NewListModel(profiles, activeName, title)
	p := tea.NewProgram(model, tea.WithAltScreen())
	p.Run()
}

// StartInfo starts an info view TUI
func StartInfo(activeProfile *persona.Profile, title string) {
	model := NewInfoModel(activeProfile, title)
	p := tea.NewProgram(model, tea.WithAltScreen())
	p.Run()
}
