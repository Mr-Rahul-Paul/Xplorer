package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

const APP = "nvim" // this app will open text files

// this file is storing the bubble tea model
// why is it confusingly named ? idk AI told me to
type Model struct {
	CurrentPath   string
	Entries       []Entry
	SelectedIndex int
	StatusMessage string
	ShowHidden    bool
	Width         int
	Height        int
}

func NewModel(path string, entries []Entry) Model {
	return Model{
		CurrentPath:   path,
		Entries:       entries,
		SelectedIndex: 0,
		StatusMessage: "",
		ShowHidden:    false,
		Width:         0,
		Height:        0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// this kinda depends on the next view
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		if msg != nil {
			m.StatusMessage = msg.Error()
		} else {
			m.StatusMessage = "Returned from editor"
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			} else {
				if len(m.Entries) != 0 {
					m.SelectedIndex = len(m.Entries) - 1
				}
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

			if selectedEntry.Type == SymlinkEntry && selectedEntry.IsBrokenSymlink {
				// dont open broken links
				m.StatusMessage = "Broken symlink"
				return m, nil
			}
			// if != DirectorDirectoryEntry ??? open the file
			if selectedEntry.Type == DirectoryEntry {
				//this is action so we can read disk it works
				entries, err := ReadDirectory(selectedEntry.FullPath, m.ShowHidden)
				//check for err
				if err != nil {
					m.StatusMessage = err.Error()
					return m, nil
				}

				m.CurrentPath = selectedEntry.FullPath // get in thre folder
				m.Entries = entries
				m.SelectedIndex = 0
				m.StatusMessage = ""
			} else {
				cmd := exec.Command(APP, selectedEntry.FullPath)
				// cmd.Start()
				// if err := cmd.Start(); err != nil {
				// 	m.StatusMessage = err.Error()
				// 	return m, nil
				// }
				// m.StatusMessage = "Opened " + selectedEntry.Name
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
					return err
				})
			}

		case "r":
			entries, err := ReadDirectory(m.CurrentPath, m.ShowHidden)
			if err != nil {
				m.StatusMessage = err.Error()
				return m, nil
			}

			m.Entries = entries
			m.SelectedIndex = 0
			m.StatusMessage = "Refreshed"

		case "backspace":
			parentPath := filepath.Dir(m.CurrentPath)

			entries, err := ReadDirectory(parentPath, m.ShowHidden)

			if err != nil {
				m.StatusMessage = err.Error()
				return m, nil
			}

			m.CurrentPath = parentPath // get in thre folder
			m.Entries = entries
			m.SelectedIndex = 0
			m.StatusMessage = ""

		case "h":
			m.StatusMessage = "Hidden files toggled"
			// toggles the state and not hard codes it
			m.ShowHidden = !m.ShowHidden

			entries, err := ReadDirectory(m.CurrentPath, m.ShowHidden)
			if err != nil {
				m.StatusMessage = err.Error()
				return m, nil
			}
			// we reload the model - readdisk is happening right (or happened above) ?
			m.Entries = entries
			m.SelectedIndex = 0

			if m.ShowHidden {
				m.StatusMessage = "showing hidden files"
			} else {
				m.StatusMessage = "not showing hidden files"
			}
			//case ends here
		}
	}
	return m, nil
}

func (m Model) View() string {
	marker := "[F] "

	switch entry.Type()

	// view := "THIS IS THE CURRENT PATH: " + m.CurrentPath + "\n\n"

	view := "Size: "
	view += fmt.Sprintf("%dx%d", m.Width, m.Height)
	view += "\n"
	view += "THIS IS THE CURRENT PATH: " + m.CurrentPath + "\n\n"

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
