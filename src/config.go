package main

import (
	"encoding/json"
	_ "time/tzdata"
	"log"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	DailyTasks []string `json:"DailyTasks"`
	TimeZone   string   `json:"TimeZone"`
	loc        *time.Location
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{TimeZone: "UTC"}

	if data, err := os.ReadFile("config.json"); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			log.Printf("Warning: failed to parse config.json: %v", err)
		}
	}

	if tasks := os.Getenv("DAILY_TASKS"); tasks != "" {
		cfg.DailyTasks = parseTasks(tasks)
	}
	if tz := os.Getenv("TIMEZONE"); tz != "" {
		cfg.TimeZone = tz
	}

	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		log.Printf("Warning: invalid timezone %q, falling back to UTC", cfg.TimeZone)
		loc = time.UTC
	}
	cfg.loc = loc

	return cfg, nil
}

func (c *AppConfig) GetToday() time.Time {
	return time.Now().In(c.loc).Truncate(24 * time.Hour)
}

func parseTasks(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}
