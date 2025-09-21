package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Config holds all pomodoro configuration
type Config struct {
	WorkDuration    time.Duration `json:"work_duration"`     // Work period duration
	BreakDuration   time.Duration `json:"break_duration"`    // Break duration after work
	WarningTime     time.Duration `json:"warning_time"`      // Warning time before sleep
	ExtendDuration  time.Duration `json:"extend_duration"`   // How long extension lasts
	AlwaysOn        bool          `json:"always_on"`         // Keep timer running continuously
	ScheduleEnabled bool          `json:"schedule_enabled"`  // Enable scheduled start times
	ScheduleStart   string        `json:"schedule_start"`    // Start time (HH:MM format)
	ScheduleEnd     string        `json:"schedule_end"`      // End time (HH:MM format)
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		WorkDuration:    50 * time.Minute,
		BreakDuration:   10 * time.Minute,
		WarningTime:     5 * time.Minute,
		ExtendDuration:  5 * time.Minute,
		AlwaysOn:        false,
		ScheduleEnabled: false,
		ScheduleStart:   "09:00",
		ScheduleEnd:     "18:00",
	}
}

// ConfigPath returns the path to the config file
func ConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "pomoduru", "config.json")
}

// LoadConfig loads configuration from file, creates default if not exists
func LoadConfig() (*Config, error) {
	configPath := ConfigPath()
	
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}
	
	// If config file doesn't exist, create with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := SaveConfig(config); err != nil {
			return nil, err
		}
		return config, nil
	}
	
	// Load existing config
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	
	return &config, nil
}

// SaveConfig saves configuration to file
func SaveConfig(config *Config) error {
	configPath := ConfigPath()
	
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
