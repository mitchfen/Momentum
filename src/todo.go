package main

import "strings"

// SubTask is a single checkable item within a TodoItem.
type SubTask struct {
	Text        string
	IsCompleted bool
}

// TodoItem represents a task or habit stack parsed from the config.
type TodoItem struct {
	Text         string
	DisplayTitle string
	SubTasks     []SubTask
}

// IsCompleted returns true when all subtasks are done.
func (t TodoItem) IsCompleted() bool {
	if len(t.SubTasks) == 0 {
		return false
	}
	for _, s := range t.SubTasks {
		if !s.IsCompleted {
			return false
		}
	}
	return true
}

// CompletedCount returns the number of completed subtasks.
func (t TodoItem) CompletedCount() int {
	n := 0
	for _, s := range t.SubTasks {
		if s.IsCompleted {
			n++
		}
	}
	return n
}

// ParseTodoItem parses a task config string into a TodoItem.
//
// Supported formats:
//
//	"Task name"                          → single task
//	"Step 1 + Step 2 + Step 3"           → unnamed habit stack
//	"My Stack: Step 1 + Step 2 + Step 3" → named habit stack
func ParseTodoItem(text string, getState func(subTaskText string) bool) TodoItem {
	var displayTitle string
	var subTaskTexts []string

	if idx := strings.Index(text, ":"); idx >= 0 {
		displayTitle = strings.TrimSpace(text[:idx])
		subTaskTexts = splitSubTasks(text[idx+1:])
	} else {
		subTaskTexts = splitSubTasks(text)
		if len(subTaskTexts) == 1 {
			displayTitle = subTaskTexts[0]
		}
	}

	subTasks := make([]SubTask, 0, len(subTaskTexts))
	for _, st := range subTaskTexts {
		subTasks = append(subTasks, SubTask{
			Text:        st,
			IsCompleted: getState(st),
		})
	}

	return TodoItem{
		Text:         text,
		DisplayTitle: displayTitle,
		SubTasks:     subTasks,
	}
}

func splitSubTasks(s string) []string {
	parts := strings.Split(s, "+")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}
