package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	form *huh.Form
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	if m.form.State == huh.StateCompleted {
		return m, tea.Quit
	}

	return m, cmd
}

func (m model) View() tea.View {
	if m.form.State == huh.StateCompleted {
		body := lipgloss.NewStyle().
			Padding(1, 2).
			Render(fmt.Sprintf(
				"✓ Form completed!\n\nFirst name: %s\nSecond name: %s",
				m.form.GetString("name1"),
				m.form.GetString("name2"),
			))
		return tea.NewView(body)
	}

	// Display form with live values below
	name1 := m.form.GetString("name1")
	name2 := m.form.GetString("name2")

	liveValues := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(1, 0).
		Render(fmt.Sprintf(
			"Live values (GetString while editing):\n  name1: %q\n  name2: %q",
			name1,
			name2,
		))

	return tea.NewView(m.form.View() + "\n" + liveValues)
}

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name1").
				Title("First name").
				Description("This is the first input field"),
			huh.NewInput().
				Key("name2").
				Title("Second name").
				Description("This is the second input field with same key structure"),
		),
	).WithWidth(60)

	p := tea.NewProgram(model{form: form})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
