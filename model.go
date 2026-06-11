package main

import tea "github.com/charmbracelet/bubbletea"

// this file is storing the bubble tea model
// why is it confusingly named ? idk AI told me to
type Model struct {
	CurrentPath   string
	Entries       []Entry
	SelectedIndex int
	StatusMessage string
}

func NewModel(path string, entries []Entry) Model {
	return Model{
		CurrentPath:   path,
		Entries:       entries,
		SelectedIndex: 0,
		StatusMessage: "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// this kinda depends on the next view
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}
		case "down", "j":
			if m.SelectedIndex < len(m.Entries)-1 {
				m.SelectedIndex++
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	view := "THIS IS THE CURRENT PATH: " + m.CurrentPath + "\n\n"

	for i, entry := range m.Entries {
		cursor := " "

		if i == m.SelectedIndex {
			cursor = "> "
		}

		view += cursor + entry.Name + "\n"
	}

	view += "\n Press q or 'ctrl+c' to quit \n"

	return view
}
