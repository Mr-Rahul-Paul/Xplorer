package main

import "time"

type EntryType int

type Entry struct {
	Name            string
	FullPath        string
	Type            EntryType
	Size            int64
	ModifiedTime    time.Time
	IsBrokenSymlink bool
}

type SortMode int

// lmao sorting by time is just a state -> int ????
const (
	SortByTime SortMode = iota
	SortByName
	SortByType
)

const (
	FileEntry EntryType = iota
	DirectoryEntry
	SymlinkEntry
	OtherEntry // not everything is a file , dir or a Symlink.
)
