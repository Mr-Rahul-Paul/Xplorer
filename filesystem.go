package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ReadDirectory(path string, showHidden bool, sortMode SortMode) ([]Entry, error) {
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

		// entryType := FileEntry /// this is wrong lil bro
		// Linux also has sockets, named pipes, and device files.
		entryType := OtherEntry
		// mode symlink checkss if symlink bit is set
		// direntry.Type() returns fs type bits
		if dirEntry.Type()&os.ModeSymlink != 0 {
			entryType = SymlinkEntry
		} else if dirEntry.IsDir() {
			entryType = DirectoryEntry
		} else if dirEntry.Type().IsRegular() {
			entryType = FileEntry
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
			Name: dirEntry.Name(),
			// no need to fetch again
			// FullPath:        filepath.Join(path, dirEntry.Name()),
			FullPath:        fullPath,
			Type:            entryType,
			Size:            info.Size(),
			ModifiedTime:    info.ModTime(),
			IsBrokenSymlink: isBrokenSymlink,
		}

		entries = append(entries, entry)
	}
	// this sorts the arr ... wierd syntax ngl
	SortEntries(entries, sortMode)

	return entries, nil
}

func SortEntries(entries []Entry, mode SortMode) {
	sort.SliceStable(entries, func(i, j int) bool {
		switch mode {
		case SortByName:
			left := strings.ToLower(entries[i].Name)
			right := strings.ToLower(entries[j].Name)
			return left < right

		case SortByType:
			if entries[i].Type == entries[j].Type {
				left := strings.ToLower(entries[i].Name)
				right := strings.ToLower(entries[j].Name)
				return left < right
			}
			return entries[i].Type < entries[j].Type

		case SortByTime:
			fallthrough
		default:
			return entries[i].ModifiedTime.After(entries[j].ModifiedTime)
		}
	})
}

func ReadNearestExisitingDirectory(path string, showHidden bool, sortMode SortMode) (string, []Entry, error) {
	currentPath := path

	for {
		entries, err := ReadDirectory(currentPath, showHidden, sortMode)
		if err == nil {
			return currentPath, entries, nil
		}

		if !os.IsNotExist(err) {
			return "", nil, err
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			return "", nil, err
		}

		currentPath = parentPath
	}
}
