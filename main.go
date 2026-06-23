package main

import (
	"log"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	startPath, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)

	}
	entries, err := ReadDirectory(startPath, false)
	//init the state of the TUI
	model := NewModel(startPath, entries)

	if err != nil {
		model.StatusMessage = err.Error()
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatal(err)

	}
}
