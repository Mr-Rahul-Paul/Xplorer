package main

import (
	"fmt"
	"os"
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

type EditorFinishedMsg struct {
	Err error
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

	case EditorFinsishedMsg:
		if msg.Err != nil {
			m.StatusMessage = "Failed to open file: " + msg.Err.Error()
		} else {
			m.StatusMessage = "returned from editor"
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
			if selectedEntry.Type == OtherEntry {
				m.StatusMessage = "Unsupported file type"
				return m, nil
			}
			// os.Stat checks the selected path immediately before using it.
			// os.Stat follows the entry to its target
			// validate before entering

			info, err := os.Stat(selectedEntry.FullPath)
			if err != nil {
				if os.IsNotExist(err) {
					entries, readErr := ReadDirectory(m.CurrentPath, m.ShowHidden)
					if readErr != nil {
						m.StatusMessage = readErr.Error()
						return m, nil
					}

					m.Entries = entries
					m.SelectedIndex = 0
					m.StatusMessage = "Path no longer exists"
					return m, nil
				}
				m.StatusMessage = err.Error()
				return m, nil
			}
			// if != DirectorDirectoryEntry ??? open the file
			if selectedEntry.Type == DirectoryEntry ||
				(selectedEntry.Type == SymlinkEntry && info.IsDir()) {
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
					return EditorFinishedMsg{Err: err}
				})
			}

		case "r":
			path, entries, err := ReadNearestExisitingDirectory(
				m.CurrentPath,
				m.ShowHidden,
			)
			if err != nil {
				m.StatusMessage = err.Error()
				return m, nil
			}

			if path != m.CurrentPath {
				m.StatusMessage = "Directory was removed , moved to nearest parent"
			} else {
				m.StatusMessage = "Refreshed"
			}

			m.CurrentPath = path
			m.Entries = entries
			m.SelectedIndex = 0

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
			// TODO : IN FUTURE ADD TOGGLE FOR MARKERS
			// case "f":
			// 	m.StatusMessage = "marker toggled"
			// 	// toggles the state and not hard codes it
			// 	m.ShowHidden = !m.ShowHidden

			// 	entries, err := ReadDirectory(m.CurrentPath, m.ShowHidden)
			// 	if err != nil {
			// 		m.StatusMessage = err.Error()
			// 		return m, nil
			// 	}
			// 	// we reload the model - readdisk is happening right (or happened above) ?
			// 	m.Entries = entries
			// 	m.SelectedIndex = 0

			// 	if m.ShowHidden {
			// 		m.StatusMessage = "showing hidden files"
			// 	} else {
			// 		m.StatusMessage = "not showing hidden files"
			// 	}
			//case ends here
		}
	}
	return m, nil
}

func (m Model) View() string {
	// view := "THIS IS THE CURRENT PATH: " + m.CurrentPath + "\n\n"

	view := "Size: "
	view += fmt.Sprintf("%dx%d", m.Width, m.Height)
	view += "\n"
	view += "THIS IS THE CURRENT PATH: " + m.CurrentPath + "\n\n"

	if len(m.Entries) == 0 {
		view += " Empty directory\n"
	}

	for i, entry := range m.Entries {

		marker := "[File] "

		switch entry.Type {
		case DirectoryEntry:
			marker = "[Directory] "
		case SymlinkEntry:
			marker = "[Link] "
		case OtherEntry:
			marker = "[Other] "

		}

		if entry.IsBrokenSymlink {
			marker = "[!broken] "
		}

		cursor := "  " // i am keeping the single space cause it looks cool
		if i == m.SelectedIndex {
			cursor = "> "
		}

		view += cursor + marker + entry.Name + "\n"
	}
	if m.StatusMessage != "" {
		view += "\n" + m.StatusMessage + "\n"
	}
	view += "\n Press q or 'ctrl+c' to quit \n"

	return view
}
