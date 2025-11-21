package tui

import (
	"fmt"
	"persona/internal/persona"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ListModel struct {
	profiles     []persona.Profile
	activeName   string
	title        string
}

func NewListModel(profiles []persona.Profile, activeName, title string) ListModel {
	return ListModel{
		profiles:   profiles,
		activeName: activeName,
		title:      title,
	}
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ListModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render(m.title))
	b.WriteString("\n\n")

	if len(m.profiles) == 0 {
		b.WriteString(MutedStyle.Render("No profiles found."))
		return BorderStyle.Render(b.String())
	}

	// Profiles
	for i, profile := range m.profiles {
		activeMarker := ""
		if profile.Name == m.activeName {
			activeMarker = " " + ActiveMarkerStyle
		}

		// Profile header
		nameStyle := ValueStyle.Bold(true)
		if profile.Name == m.activeName {
			nameStyle = nameStyle.Foreground(SuccessColor)
		}
		b.WriteString(fmt.Sprintf("%d. %s%s\n", i+1, nameStyle.Render(profile.Name), activeMarker))

		// Profile details
		detailStyle := MutedStyle.MarginLeft(3)
		b.WriteString(detailStyle.Render(fmt.Sprintf("User: %s\n", profile.User)))
		b.WriteString(detailStyle.Render(fmt.Sprintf("Email: %s\n", profile.Email)))

		if profile.SigningKey != "" {
			b.WriteString(detailStyle.Render(fmt.Sprintf("Signing Key: %s\n", profile.SigningKey)))
		}

		if profile.Description != "" {
			b.WriteString(detailStyle.Render(fmt.Sprintf("Description: %s\n", profile.Description)))
		}

		if i < len(m.profiles)-1 {
			b.WriteString("\n")
		}
	}

	// Footer
	b.WriteString("\n")
	footer := fmt.Sprintf("Total: %d profile(s)", len(m.profiles))
	b.WriteString(MutedStyle.Render(footer))

	// Help
	b.WriteString("\n\n")
	help := "Press Esc or 'q' to exit"
	b.WriteString(HelpStyle.Render(help))

	return BorderStyle.Render(b.String())
}

