package main

import (
	"fmt"
	"os"

	"github.com/aniketvish/pomoduru/internal/config"
	"github.com/aniketvish/pomoduru/internal/timer"
	"github.com/aniketvish/pomoduru/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Create timer
	t := timer.NewTimer(cfg)

	// Create and start scheduler if enabled
	scheduler := timer.NewScheduler(cfg, t)
	scheduler.Start()
	defer scheduler.Stop()

	// Create UI model
	model := ui.NewModel(cfg, t)

	// Start the TUI
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
