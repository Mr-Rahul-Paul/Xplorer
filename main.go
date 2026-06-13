package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Starting point
	entries, err := ReadDirectory(".", false)

	if err != nil {
		log.Fatal(err)
	}
	//init the state of the TUI
	model := NewModel(".", entries)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		log.Fatal(err)

	}
}
