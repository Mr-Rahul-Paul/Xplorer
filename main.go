package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Starting point
	entries, err := ReadDirectory(".", false)
	//init the state of the TUI
	model := NewModel(".", entries)

	if err != nil {
		model.StatusMessage = err.Error()
	}

	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatal(err)

	}
}
