package tui

import (
	"fmt"
	"os/exec"
	"persona/internal/persona"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type InfoModel struct {
	activeProfile *persona.Profile
	gitUserName   string
	gitUserEmail  string
	title         string
}

func NewInfoModel(activeProfile *persona.Profile, title string) InfoModel {
	model := InfoModel{
		activeProfile: activeProfile,
		title:         title,
	}

	// Get git config
	out, err := exec.Command("git", "config", "user.name").Output()
	model.gitUserName = strings.TrimSpace(string(out))
	if err != nil || model.gitUserName == "" {
		model.gitUserName = "Not configured"
	}

	out, err = exec.Command("git", "config", "user.email").Output()
	model.gitUserEmail = strings.TrimSpace(string(out))
	if err != nil || model.gitUserEmail == "" {
		model.gitUserEmail = "Not configured"
	}

	return model
}

func (m InfoModel) Init() tea.Cmd {
	return nil
}

func (m InfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m InfoModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(TitleStyle.Render(m.title))
	b.WriteString("\n\n")

	// Active Profile Section
	if m.activeProfile != nil {
		b.WriteString(HeaderStyle.Render("Active Profile"))
		b.WriteString("\n")
		nameStyle := ValueStyle.Bold(true)
		b.WriteString(fmt.Sprintf("  %s %s\n", ActiveMarkerStyle, nameStyle.Render(m.activeProfile.Name)))
		
		detailStyle := MutedStyle.MarginLeft(3)
		b.WriteString(detailStyle.Render(fmt.Sprintf("User: %s\n", m.activeProfile.User)))
		b.WriteString(detailStyle.Render(fmt.Sprintf("Email: %s\n", m.activeProfile.Email)))
		
		if m.activeProfile.SigningKey != "" {
			b.WriteString(detailStyle.Render(fmt.Sprintf("Signing Key: %s\n", m.activeProfile.SigningKey)))
		}
		
		if m.activeProfile.Description != "" {
			b.WriteString(detailStyle.Render(fmt.Sprintf("Description: %s\n", m.activeProfile.Description)))
		}
		b.WriteString("\n")
	} else {
		b.WriteString(ErrorStyle.Render("❌ No active profile set."))
		b.WriteString("\n")
		b.WriteString(MutedStyle.Render("Use `persona add` to add a profile or `persona use` to select a profile."))
		b.WriteString("\n\n")
	}

	// Git Identity Section
	b.WriteString(HeaderStyle.Render("Git Identity"))
	b.WriteString("\n")
	detailStyle := MutedStyle.MarginLeft(3)
	b.WriteString(detailStyle.Render(fmt.Sprintf("User Name: %s\n", m.gitUserName)))
	b.WriteString(detailStyle.Render(fmt.Sprintf("Email: %s\n", m.gitUserEmail)))

	// Help
	b.WriteString("\n")
	help := "Press Esc or 'q' to exit"
	b.WriteString(HelpStyle.Render(help))

	return BorderStyle.Render(b.String())
}

