package main

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

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
			} else {
				m.SelectedIndex = len(m.Entries) - 1
			}
		case "down", "j":
			if m.SelectedIndex < len(m.Entries)-1 {
				m.SelectedIndex++
			} else {
				m.SelectedIndex = 0
			}
		case "enter":
			if len(m.Entries) == 0 {
				return m, nil
			}

			selectedEntry := m.Entries[m.SelectedIndex]
			// if != DirectorDirectoryEntry ??? open the file
			if selectedEntry.Type == DirectoryEntry {
				//this is action so we can read disk it works
				entries, err := ReadDirectory(selectedEntry.FullPath)
				//check for err
				if err != nil {
					m.StatusMessage = err.Error()
					return m, nil
				}

				m.CurrentPath = selectedEntry.FullPath // get in thre folder
				m.Entries = entries
				m.SelectedIndex = 0
				m.StatusMessage = ""
			}

		case "backspace":
			parentPath := filepath.Dir(m.CurrentPath)

			entries, err := ReadDirectory(parentPath)

			if err != nil {
				m.StatusMessage = err.Error()
				return m, nil
			}

			m.CurrentPath = parentPath // get in thre folder
			m.Entries = entries
			m.SelectedIndex = 0
			m.StatusMessage = ""
		}
	}
	return m, nil
}

func (m Model) View() string {
	view := "THIS IS THE CURRENT PATH: " + m.CurrentPath + "\n\n"

	for i, entry := range m.Entries {
		cursor := " " // i am keeping the single space cause it looks cool

		if i == m.SelectedIndex {
			cursor = "> "
		}

		//show err

		view += cursor + entry.Name + "\n"
	}
	if m.StatusMessage != "" {
		view += "\n" + m.StatusMessage + "\n"
	}
	view += "\n Press q or 'ctrl+c' to quit \n"

	return view
}
