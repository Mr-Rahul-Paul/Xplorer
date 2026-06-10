package main

import (
	"os"
	"path/filepath"
	"sort"
)

func ReadDirectory(path string) ([]Entry, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var entries []Entry

	for _, dirEntry := range dirEntries {
		entryType := FileEntry

		if dirEntry.IsDir() {
			entryType = DirectoryEntry
		}

		info, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		entry := Entry{
			Name:         dirEntry.Name(),
			FullPath:     filepath.Join(path, dirEntry.Name()),
			Type:         entryType,
			ModifiedTime: info.ModTime(),
		}

		entries = append(entries, entry)
	}
	// this sorts the arr ... wierd syntax ngl
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ModifiedTime.After(entries[j].ModifiedTime)
	})
	return entries, nil
}
