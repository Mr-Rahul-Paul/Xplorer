# Go TUI File Explorer — Design Doc

## 1. Goal

Build a terminal file explorer in Go.

The main goal is to learn:

- Go project structure
- filesystem handling
- Bubble Tea TUI architecture
- state-based app design
- error handling in real filesystem cases

The app should let the user:

- browse directories
- move selection up/down
- enter folders
- go back to parent folder
- open files using Kate
- handle filesystem errors without crashing

---

## 2. Core Mental Model

This is not a recursive scanner.

The app always shows one current directory.

When the user enters a folder, the app reads only that folder and replaces the current view.

Main rule:

Movement does not read disk.  
Actions read disk.

Examples:

- Up / Down: only changes selected index
- Enter: validates selected path
- Backspace: validates parent path
- r: reloads current directory

The app is basically:

State + Filesystem + Input + Render

---

## 3. Main Architecture

The app has 4 main parts:

1. State
2. Filesystem
3. Input
4. UI / Render

### Responsibility Split

State:

- stores current app data
- remembers current path
- remembers selected item
- stores current entries
- stores status/error message

Filesystem:

- reads directories
- classifies entries
- checks if paths exist
- detects broken symlinks
- returns errors safely

Input:

- handles keypresses
- decides what action to perform
- updates state
- calls filesystem only when needed

UI / Render:

- displays current state
- shows files/folders
- highlights selected entry
- shows status/error message

---

## 4. State Design

The app state should store only what the UI and input system need.

Suggested Bubble Tea model:

```go
type Model struct {
    CurrentPath   string
    Entries       []Entry
    SelectedIndex int
    StatusMessage string
    SortMode      SortMode
    ShowHidden    bool
    Width         int
    Height        int
}
```

```go
type Entry struct {
    Name            string
    FullPath        string
    Type            EntryType
    Size            int64
    ModifiedTime    time.Time
    IsBrokenSymlink bool
}
```

```go
type EntryType int

const (
    FileEntry EntryType = iota
    DirectoryEntry
    SymlinkEntry
    OtherEntry
)
```

---

## 5. Filesystem Layer

The filesystem layer should only deal with disk/filesystem logic.

Responsibilities:

read the directory
skip hidden files if ShowHidden is false
classify entries
detect symlinks
detect broken symlinks
get size
get modified time
sort entries
return clean errors

This layer should not:

render UI
read keyboard input
change selected index
open Kate directly

## Input Rules

Controls:

Up / k: move selection up
Down / j: move selection down
Enter: open selected item
Backspace: go to parent directory
r: refresh current directory
h: toggle hidden files
s: cycle sort mode later
q / ctrl+c: quit

Important rule:

Movement keys should not read disk.

Only action keys should validate/read disk.

No reload for:

Up
Down
scrolling
moving selection

Reload/validate for:

Enter
Backspace
r
h
opening file
opening folder

--- 
## Sorting Rules

Default sorting:

Newest modified first.

This means recently changed files/folders appear at the top.

Suggested enum:
```go
type SortMode int

const (
    SortByTime SortMode = iota
    SortByName
    SortByType
)
```

11. Error Rules
Permission Denied

Case:

User tries to enter a directory they cannot read.

Behavior:

stay in current directory
show message
do not crash

Message:

Permission denied: /path/name
Path No Longer Exists

Case:

Selected file/folder was deleted before opening.

Behavior:

show message
refresh current directory
keep user in current view
do not crash

Message:

Path no longer exists
File Deleted Before Opening

Case:

User presses Enter on selected file, but another process deleted it.

Behavior:

validate path exists before opening
if missing, show message
refresh current directory

Message:

File was deleted
Current Directory Deleted While Browsing

Case:

User is inside:

/home/rahul/project/temp

Another terminal deletes:

temp

Behavior:

on refresh/navigation, detect current path is missing
move to nearest existing parent
show message

Message:

Directory was removed, moved to parent

Example:

If this is deleted:

/home/rahul/project/temp

Move to:

/home/rahul/project

If that is also deleted, keep moving upward until an existing parent is found.

Broken Symlink

Case:

Symlink points to a missing target.

Example:

shortcut -> /deleted/path

Behavior:

show it as a symlink
mark it as broken
on Enter, do not crash
show message

Message:

Broken symlink
Kate/Open Command Fails

Case:

Kate is not installed or command fails.

Behavior:

stay in current directory
show message
do not crash

Message:

Failed to open file
Empty Directory

Case:

Directory has no visible entries.

Behavior:

show empty directory message
selectedIndex stays 0
app does not crash

Message:

Empty directory
12. Opening Rules

For v1:

Enter on folder: navigate into folder
Enter on file: open with Kate
Enter on broken symlink: show "Broken symlink"
Enter on other type: show "Unsupported file type"

Opening command:

kate /path/to/file

In Go:

exec.Command("kate", path).Start()

Use Start(), not Run().

Why:

Start opens Kate without blocking the TUI forever.

Run waits until Kate exits, which is not ideal for a TUI.

13. Folder Navigation Rules

When Enter is pressed on a directory:

check path still exists
check it is readable
read directory entries
update CurrentPath
update Entries
reset SelectedIndex to 0
clear or update StatusMessage

If reading fails:

stay in current path
show error message
do not crash
14. File Opening Rules

When Enter is pressed on a file:

check path still exists
if missing, show "File was deleted"
refresh current directory
if exists, open with Kate
if Kate fails, show "Failed to open file"

Opening file should not change CurrentPath.

15. Symlink Rules

For v1:

show symlink as symlink
detect if target is missing
if target missing, mark broken
pressing Enter on broken symlink shows error

Simple behavior:

Broken symlink: show error
Normal symlink: treat carefully later

Do not overbuild symlink handling in v1.

16. Hidden Files

Linux hidden files start with dot.

Examples:

.git
.env
.config

Default:

ShowHidden = false

Key:

h

Behavior:

toggle ShowHidden
reload current directory
clamp selected index