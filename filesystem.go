package main

import (
	"os"
	"path/filepath"
	"sort"
)

func ReadDirectory(path string, showHidden bool) ([]Entry, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var entries []Entry

	for _, dirEntry := range dirEntries {
		//what does it mean - we are checking of files name starts with
		// "." we continue or not , so dot files and folders are hidden
		// G--DAMn
		if !showHidden && dirEntry.Name()[0] == '.' {
			continue
		}

		entryType := FileEntry
		// mode symlink checkss if symlink bit is set
		// direntry.Type() returns fs type bits
		if dirEntry.Type()&os.ModeSymlink != 0 {
			entryType = SymlinkEntry
		} else if dirEntry.IsDir() {
			entryType = DirectoryEntry
		}

		//after entry identified
		fullPath := filepath.Join(path, dirEntry.Name())
		isBrokenSymlink := false

		if entryType == SymlinkEntry {
			// os.Stat follows the symlink to its target
			if _, err := os.Stat(fullPath); err != nil && os.IsNotExist(err) {
				isBrokenSymlink = true
			}
		}

		info, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		entry := Entry{
			Name:            dirEntry.Name(),
			FullPath:        filepath.Join(path, dirEntry.Name()),
			Type:            entryType,
			Size:            info.Size(),
			ModifiedTime:    info.ModTime(),
			IsBrokenSymlink: isBrokenSymlink,
		}

		entries = append(entries, entry)
	}
	// this sorts the arr ... wierd syntax ngl
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ModifiedTime.After(entries[j].ModifiedTime)
	})
	return entries, nil
}
