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
