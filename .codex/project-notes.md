# Xplorer Codex Notes

## Reply Size

- Keep replies minimal and token-saving.
- Use short explanations.
- Avoid long paragraphs and large code dumps.

## Aim

- Build a Go Bubble Tea terminal file explorer step by step.
- The user is learning Go, so teach through small compiling steps.
- Do not add features before the requested build step.

## Teaching Style

- Give one step at a time.
- Explain only the Go concepts that appear in that step.
- Prefer beginner-readable Go over clever code.
- Let the user write the code unless they explicitly ask Codex to edit files.
- Check the user's code before moving forward when asked.

## Required Step Format

1. What we are building
2. Why this file/code exists
3. Code block
4. What this code does in 3-5 bullets
5. Command to run/test
6. Stop and wait for me

## Architecture

- State stores current path, entries, selected index, and status message.
- Filesystem reads directories and returns entries.
- Input handles keypresses and updates state.
- UI renders the current state.
- Movement does not read disk.
- Actions read disk.

## Current Build Order

1. Create Go module and basic main.go. Done.
2. Create Entry and EntryType. Done.
3. Read current directory and print entries. Done.
4. Sort entries by modified time, newest first. Done.
5. Add Bubble Tea model. Current/next.

## UI Direction

- Target layout: left sidebar plus main content area.
- Left sidebar sections: major folders, saved/recent/important paths, disk size filled/remaining.
- Main area: top path bar, large content panel below.
- Keep this design in mind when rendering the TUI later.
