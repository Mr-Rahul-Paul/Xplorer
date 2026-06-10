package main

type EntryType int

type Entry struct {
	Name     string
	FullPath string
	Type     EntryType
}

const (
	FileEntry EntryType = iota
	DirectoryEntry
	SymlinkEntry
	OtherEntry
)
