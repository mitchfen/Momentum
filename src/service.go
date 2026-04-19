package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TodoService persists daily subtask completion state to a JSON file.
// State is keyed by date (yyyy-MM-dd) and cleared on the next day.
type TodoService struct {
	cfg      *AppConfig
	filePath string
	mu       sync.Mutex
	// state: date → taskText → subTaskText → isCompleted
	state map[string]map[string]map[string]bool
}

func NewTodoService(cfg *AppConfig) *TodoService {
	svc := &TodoService{
		cfg:      cfg,
		filePath: "data/todo-state.json",
		state:    make(map[string]map[string]map[string]bool),
	}
	svc.loadFromFile()
	return svc
}

func (s *TodoService) GetSubTaskState(date time.Time, taskText, subTaskText string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	dateKey := date.Format("2006-01-02")
	if tasks, ok := s.state[dateKey]; ok {
		if subTasks, ok := tasks[taskText]; ok {
			return subTasks[subTaskText]
		}
	}
	return false
}

func (s *TodoService) SetSubTaskState(date time.Time, taskText, subTaskText string, isCompleted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dateKey := date.Format("2006-01-02")
	if _, ok := s.state[dateKey]; !ok {
		s.state[dateKey] = make(map[string]map[string]bool)
	}
	if _, ok := s.state[dateKey][taskText]; !ok {
		s.state[dateKey][taskText] = make(map[string]bool)
	}
	s.state[dateKey][taskText][subTaskText] = isCompleted
	s.cleanupOldDates()
	s.saveToFile()
}

func (s *TodoService) cleanupOldDates() {
	today := s.cfg.GetToday().Format("2006-01-02")
	for key := range s.state {
		if key != today {
			delete(s.state, key)
		}
	}
}

func (s *TodoService) loadFromFile() {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}
	var state map[string]map[string]map[string]bool
	if err := json.Unmarshal(data, &state); err != nil {
		log.Printf("Error loading state: %v", err)
		return
	}
	s.state = state
	s.cleanupOldDates()
}

func (s *TodoService) saveToFile() {
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		log.Printf("Error creating data directory: %v", err)
		return
	}
	data, err := json.Marshal(s.state)
	if err != nil {
		log.Printf("Error marshaling state: %v", err)
		return
	}
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		log.Printf("Error saving state: %v", err)
	}
}
