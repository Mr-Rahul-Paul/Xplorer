package main

// this file is storing the bubble tea model
// why is it confusingly named ? idk AI told me to
type Model struct {
	CurrentPath   string
	Entires       []Entry
	SelectedIndex int
	StatusMessage string
}

func NewModel(path string, entries []Entry) Model {
	return Model{
		CurrentPath:   path,
		Entires:       entries,
		SelectedIndex: 0,
		StatusMessage: "",
	}
}
