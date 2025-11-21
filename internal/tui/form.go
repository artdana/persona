package tui

import (
	"fmt"
	"persona/internal/persona"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormField struct {
	Label       string
	Value       string
	Placeholder string
	Required    bool
	Input       textinput.Model
}

type FormModel struct {
	fields      []FormField
	currentIdx  int
	title       string
	cancelled   bool
	onSubmit    func(map[string]string) bool
	errors      map[int]string
}

func NewFormModel(title string, fields []FormField, onSubmit func(map[string]string) bool) FormModel {
	model := FormModel{
		fields:    make([]FormField, len(fields)),
		currentIdx: 0,
		title:     title,
		cancelled: false,
		onSubmit:  onSubmit,
		errors:    make(map[int]string),
	}

	// Initialize text inputs
	for i := range fields {
		ti := textinput.New()
		ti.Placeholder = fields[i].Placeholder
		ti.SetValue(fields[i].Value)
		ti.CharLimit = 200
		ti.Width = 50
		if i == 0 {
			ti.Focus()
		}
		model.fields[i] = FormField{
			Label:       fields[i].Label,
			Value:       fields[i].Value,
			Placeholder: fields[i].Placeholder,
			Required:    fields[i].Required,
			Input:       ti,
		}
	}

	return model
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			if m.currentIdx < len(m.fields)-1 {
				// Move to next field
				m.fields[m.currentIdx].Input.Blur()
				m.currentIdx++
				m.fields[m.currentIdx].Input.Focus()
				return m, textinput.Blink
			} else {
				// Submit form
				if m.validate() {
					if m.onSubmit != nil {
						values := m.getValues()
						if m.onSubmit(values) {
							return m, tea.Quit
						}
					} else {
						return m, tea.Quit
					}
				}
			}
		case "shift+tab", "up":
			if m.currentIdx > 0 {
				m.fields[m.currentIdx].Input.Blur()
				m.currentIdx--
				m.fields[m.currentIdx].Input.Focus()
				return m, textinput.Blink
			}
		case "tab", "down":
			if m.currentIdx < len(m.fields)-1 {
				m.fields[m.currentIdx].Input.Blur()
				m.currentIdx++
				m.fields[m.currentIdx].Input.Focus()
				return m, textinput.Blink
			}
		}
	}

	// Update focused input
	var cmd tea.Cmd
	m.fields[m.currentIdx].Input, cmd = m.fields[m.currentIdx].Input.Update(msg)
	return m, cmd
}

func (m FormModel) validate() bool {
	m.errors = make(map[int]string)
	valid := true

	for i, field := range m.fields {
		value := strings.TrimSpace(field.Input.Value())
		if field.Required && value == "" {
			m.errors[i] = fmt.Sprintf("%s is required", field.Label)
			valid = false
		}
	}

	return valid
}

func (m FormModel) getValues() map[string]string {
	values := make(map[string]string)
	for _, field := range m.fields {
		value := strings.TrimSpace(field.Input.Value())
		// For required fields, use original value if empty (for editing scenarios)
		// For optional fields, allow empty values
		if value == "" && field.Required && field.Value != "" {
			value = field.Value
		}
		values[field.Label] = value
	}
	return values
}

func (m FormModel) View() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(m.title))
	b.WriteString("\n\n")

	maxLabelWidth := 0
	for _, field := range m.fields {
		labelText := field.Label
		if field.Required {
			labelText += " *"
		}
		labelText += ":"
		if len(labelText) > maxLabelWidth {
			maxLabelWidth = len(labelText)
		}
	}

	inputContentWidth := 50
	borderedStyle := lipgloss.NewStyle().
		Width(inputContentWidth).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(FocusedColor).
		Padding(0, 1)
	testBordered := borderedStyle.Render("")
	actualBorderedWidth := lipgloss.Width(testBordered)

	for i, field := range m.fields {
		labelText := field.Label
		if field.Required {
			labelText += " *"
		}
		labelText += ":"
		
		label := LabelStyle.Width(maxLabelWidth).Align(lipgloss.Right).Render(labelText)

		input := field.Input.View()
		
		if i == m.currentIdx {
			input = borderedStyle.Render(input)
		} else {
			unfocusedStyle := lipgloss.NewStyle().Width(actualBorderedWidth)
			input = unfocusedStyle.Render(input)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Left, label, " ", input)
		b.WriteString(row)
		b.WriteString("\n")

		if err, ok := m.errors[i]; ok {
			errorPadding := strings.Repeat(" ", maxLabelWidth+1)
			b.WriteString(ErrorStyle.Render(errorPadding + err))
			b.WriteString("\n")
		}

		b.WriteString("\n")
	}

	help := "Tab/↓: Next field  Shift+Tab/↑: Previous field  Enter: Submit  Esc: Cancel"
	b.WriteString(HelpStyle.Render(help))

	return BorderStyle.Render(b.String())
}

func (m FormModel) Cancelled() bool {
	return m.cancelled
}

func (m FormModel) Values() map[string]string {
	return m.getValues()
}

func CreateProfileFormFields(profile *persona.Profile, isEdit bool) []FormField {
	fields := []FormField{
		{
			Label:       "Profile name",
			Value:       "",
			Placeholder: "",
			Required:    true,
		},
		{
			Label:       "Git user name",
			Value:       "",
			Placeholder: "",
			Required:    true,
		},
		{
			Label:       "Git email",
			Value:       "",
			Placeholder: "",
			Required:    true,
		},
		{
			Label:       "Signing key",
			Value:       "",
			Placeholder: "(optional)",
			Required:    false,
		},
		{
			Label:       "Description",
			Value:       "",
			Placeholder: "(optional)",
			Required:    false,
		},
	}

	if isEdit && profile != nil {
		fields[0].Value = profile.Name
		fields[0].Placeholder = profile.Name
		fields[1].Value = profile.User
		fields[1].Placeholder = profile.User
		fields[2].Value = profile.Email
		fields[2].Placeholder = profile.Email
		fields[3].Value = profile.SigningKey
		fields[3].Placeholder = profile.SigningKey
		fields[4].Value = profile.Description
		fields[4].Placeholder = profile.Description
	}

	return fields
}

