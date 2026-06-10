package main

import "time"

type EntryType int

type Entry struct {
	Name         string
	FullPath     string
	Type         EntryType
	ModifiedTime time.Time
}

const (
	FileEntry EntryType = iota
	DirectoryEntry
	SymlinkEntry
	OtherEntry
)
