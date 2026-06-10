package main

import (
	"os"
	"path/filepath"
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

		entry := Entry{
			Name:     dirEntry.Name(),
			FullPath: filepath.Join(path, dirEntry.Name()),
			Type:     entryType,
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
