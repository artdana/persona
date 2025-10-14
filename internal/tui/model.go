package tui

import (
	"fmt"
	"persona/internal/persona"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styles
var (
	defaultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	activeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	activeMarkerStyle  = "✅"
	cursorStyle  = "\t➜"
)


type Model struct {
	profiles      []persona.Profile
	activeName    string
	cursor 		  int
	selectedIndex int
	filter        string
}


func (m Model) Init() tea.Cmd { return nil }


func NewProfileModel(profiles []persona.Profile, activeName string) Model {
	return Model{
		profiles:      profiles,
		activeName:    activeName,
		cursor:        0,
		selectedIndex: -1,
		filter: 	   "",
	}
}


func (m Model) SelectedProfile() *persona.Profile {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.profiles) {
		return &m.profiles[m.selectedIndex]
	}

	return nil
}


// updates based on keyboard control
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.filteredProfiles())-1 {
				m.cursor++
			}
		case "enter":
			if len(m.filteredProfiles()) > 0 {
				m.selectedIndex = m.getOriginalIndex(m.cursor)
			}
			return m, tea.Quit
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.cursor = 0
			} else if msg.String() == "backspace"  || msg.String() == "delete" {
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
					m.cursor = 0
				}
			}
		}
	}
	return m, nil
}


// filters profiles based on the string
func (m Model) filteredProfiles() []persona.Profile {
	var f []persona.Profile
	
	// First, add the active profile if it matches the filter
	for _, p := range m.profiles {
		if p.Name == m.activeName {
			if m.filter == "" || strings.Contains(strings.ToLower(p.Name), strings.ToLower(m.filter)) ||
				strings.Contains(strings.ToLower(p.User), strings.ToLower(m.filter)) ||
				strings.Contains(strings.ToLower(p.Email), strings.ToLower(m.filter)) {
				f = append(f, p)
			}
			break
		}
	}
	
	// Then add all other profiles that match the filter
	for _, p := range m.profiles {
		if p.Name != m.activeName {
			if m.filter == "" || strings.Contains(strings.ToLower(p.Name), strings.ToLower(m.filter)) ||
				strings.Contains(strings.ToLower(p.User), strings.ToLower(m.filter)) ||
				strings.Contains(strings.ToLower(p.Email), strings.ToLower(m.filter)) {
				f = append(f, p)
			}
		}
	}

	return f
}


// maps the filtered index back to the original index
func (m Model) getOriginalIndex(index int) int {
	f := m.filteredProfiles()
	for i, p := range m.profiles {
		if p == f[index] {
			return i
		}
	}
	return -1
}


func (m Model) View() string {
	f := m.filteredProfiles()
	var s strings.Builder

	s.WriteString(defaultStyle.Render(fmt.Sprintf("\nSearch:%s", m.filter)))
	s.WriteString("\n\n")

	for i, p := range f {
		cursor := ""
		if i == m.cursor {
			cursor = cursorStyle
		}

		activeMarker := ""
		if p.Name == m.activeName {
			activeMarker = activeMarkerStyle
		}

		line := fmt.Sprintf("%s[%s] %s: %s", cursor, activeMarker, p.Name, p.User)

		if m.cursor == i {
			if p.Name == m.activeName {
				line = activeStyle.Render(line)
			} else {
				line = focusedStyle.Render(line)
			}
			preview := fmt.Sprintf("\t\nEmail: %s\nDescription: %s\n", p.Email, p.Description)
			line += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(preview)
		}

		s.WriteString(line + "\n")
	}

	s.WriteString("\nUse ↑/↓ to navigate, Enter to select, type to filter, escape to quit.")
	return s.String()
}
