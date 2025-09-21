package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aniketvish/pomoduru/internal/config"
)

var setCmd = flag.NewFlagSet("set", flag.ExitOnError)

func main() {
	// Subcommands
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	
	// Set flags
	setWorkDuration := setCmd.Int("work", 0, "Work duration in minutes")
	setBreakDuration := setCmd.Int("break", 0, "Break duration in minutes")
	setWarningTime := setCmd.Int("warning", 0, "Warning time in minutes before suspend")
	setExtendDuration := setCmd.Int("extend", 0, "Extension duration in minutes")
	setAlwaysOn := setCmd.Bool("always-on", false, "Enable always-on mode")
	setScheduleEnabled := setCmd.Bool("schedule-enabled", false, "Enable scheduled start times")
	setScheduleStart := setCmd.String("schedule-start", "", "Schedule start time (HH:MM)")
	setScheduleEnd := setCmd.String("schedule-end", "", "Schedule end time (HH:MM)")
	
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	
	switch os.Args[1] {
	case "show":
		showCmd.Parse(os.Args[2:])
		showConfig()
	case "set":
		setCmd.Parse(os.Args[2:])
		setConfig(setWorkDuration, setBreakDuration, setWarningTime, setExtendDuration, 
			setAlwaysOn, setScheduleEnabled, setScheduleStart, setScheduleEnd)
	case "interactive":
		interactiveConfig()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Pomoduru Configuration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  pomoduru-config show                          - Show current configuration")
	fmt.Println("  pomoduru-config set [flags]                   - Set configuration values")
	fmt.Println("  pomoduru-config interactive                   - Interactive configuration")
	fmt.Println()
	fmt.Println("Set flags:")
	fmt.Println("  --work int           Work duration in minutes (default 50)")
	fmt.Println("  --break int          Break duration in minutes (default 10)")
	fmt.Println("  --warning int        Warning time in minutes (default 5)")
	fmt.Println("  --extend int         Extension duration in minutes (default 5)")
	fmt.Println("  --always-on          Enable always-on mode")
	fmt.Println("  --schedule-enabled   Enable scheduled start times")
	fmt.Println("  --schedule-start     Schedule start time (HH:MM)")
	fmt.Println("  --schedule-end       Schedule end time (HH:MM)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  pomoduru-config set --work 45 --break 15")
	fmt.Println("  pomoduru-config set --schedule-enabled --schedule-start 09:00 --schedule-end 18:00")
}

func showConfig() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("🍅 Pomoduru Configuration")
	fmt.Println("═══════════════════════════")
	fmt.Printf("Work Duration:    %d minutes\n", int(cfg.WorkDuration.Minutes()))
	fmt.Printf("Break Duration:   %d minutes\n", int(cfg.BreakDuration.Minutes()))
	fmt.Printf("Warning Time:     %d minutes\n", int(cfg.WarningTime.Minutes()))
	fmt.Printf("Extend Duration:  %d minutes\n", int(cfg.ExtendDuration.Minutes()))
	fmt.Printf("Always On:        %t\n", cfg.AlwaysOn)
	fmt.Printf("Schedule Enabled: %t\n", cfg.ScheduleEnabled)
	fmt.Printf("Schedule Start:   %s\n", cfg.ScheduleStart)
	fmt.Printf("Schedule End:     %s\n", cfg.ScheduleEnd)
	fmt.Printf("\nConfig file: %s\n", config.ConfigPath())
}

func setConfig(work, break_, warning, extend *int, alwaysOn, scheduleEnabled *bool, 
	scheduleStart, scheduleEnd *string) {
	
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	
	changed := false
	
	if *work > 0 {
		cfg.WorkDuration = time.Duration(*work) * time.Minute
		changed = true
		fmt.Printf("Work duration set to %d minutes\n", *work)
	}
	
	if *break_ > 0 {
		cfg.BreakDuration = time.Duration(*break_) * time.Minute
		changed = true
		fmt.Printf("Break duration set to %d minutes\n", *break_)
	}
	
	if *warning > 0 {
		cfg.WarningTime = time.Duration(*warning) * time.Minute
		changed = true
		fmt.Printf("Warning time set to %d minutes\n", *warning)
	}
	
	if *extend > 0 {
		cfg.ExtendDuration = time.Duration(*extend) * time.Minute
		changed = true
		fmt.Printf("Extend duration set to %d minutes\n", *extend)
	}
	
	if *alwaysOn {
		cfg.AlwaysOn = true
		changed = true
		fmt.Printf("Always-on mode enabled\n")
	}
	
	if *scheduleEnabled {
		cfg.ScheduleEnabled = true
		changed = true
		fmt.Printf("Schedule enabled\n")
	}
	
	if *scheduleStart != "" {
		cfg.ScheduleStart = *scheduleStart
		changed = true
		fmt.Printf("Schedule start set to %s\n", *scheduleStart)
	}
	
	if *scheduleEnd != "" {
		cfg.ScheduleEnd = *scheduleEnd
		changed = true
		fmt.Printf("Schedule end set to %s\n", *scheduleEnd)
	}
	
	if changed {
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Configuration saved successfully!")
	} else {
		fmt.Println("No changes made. Use flags to set configuration values.")
	}
}

func interactiveConfig() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("🍅 Pomoduru Interactive Configuration")
	fmt.Println("══════════════════════════════════════")
	fmt.Println("Press Enter to keep current values")
	fmt.Println()
	
	// Work duration
	fmt.Printf("Work duration (minutes) [%d]: ", int(cfg.WorkDuration.Minutes()))
	if scanner.Scan() {
		if val := strings.TrimSpace(scanner.Text()); val != "" {
			if minutes, err := strconv.Atoi(val); err == nil && minutes > 0 {
				cfg.WorkDuration = time.Duration(minutes) * time.Minute
			}
		}
	}
	
	// Break duration
	fmt.Printf("Break duration (minutes) [%d]: ", int(cfg.BreakDuration.Minutes()))
	if scanner.Scan() {
		if val := strings.TrimSpace(scanner.Text()); val != "" {
			if minutes, err := strconv.Atoi(val); err == nil && minutes > 0 {
				cfg.BreakDuration = time.Duration(minutes) * time.Minute
			}
		}
	}
	
	// Warning time
	fmt.Printf("Warning time (minutes) [%d]: ", int(cfg.WarningTime.Minutes()))
	if scanner.Scan() {
		if val := strings.TrimSpace(scanner.Text()); val != "" {
			if minutes, err := strconv.Atoi(val); err == nil && minutes > 0 {
				cfg.WarningTime = time.Duration(minutes) * time.Minute
			}
		}
	}
	
	// Extend duration
	fmt.Printf("Extend duration (minutes) [%d]: ", int(cfg.ExtendDuration.Minutes()))
	if scanner.Scan() {
		if val := strings.TrimSpace(scanner.Text()); val != "" {
			if minutes, err := strconv.Atoi(val); err == nil && minutes > 0 {
				cfg.ExtendDuration = time.Duration(minutes) * time.Minute
			}
		}
	}
	
	// Always on
	fmt.Printf("Always on mode (y/n) [%t]: ", cfg.AlwaysOn)
	if scanner.Scan() {
		if val := strings.ToLower(strings.TrimSpace(scanner.Text())); val != "" {
			cfg.AlwaysOn = val == "y" || val == "yes"
		}
	}
	
	// Schedule enabled
	fmt.Printf("Schedule enabled (y/n) [%t]: ", cfg.ScheduleEnabled)
	if scanner.Scan() {
		if val := strings.ToLower(strings.TrimSpace(scanner.Text())); val != "" {
			cfg.ScheduleEnabled = val == "y" || val == "yes"
		}
	}
	
	if cfg.ScheduleEnabled {
		// Schedule start
		fmt.Printf("Schedule start time (HH:MM) [%s]: ", cfg.ScheduleStart)
		if scanner.Scan() {
			if val := strings.TrimSpace(scanner.Text()); val != "" {
				cfg.ScheduleStart = val
			}
		}
		
		// Schedule end
		fmt.Printf("Schedule end time (HH:MM) [%s]: ", cfg.ScheduleEnd)
		if scanner.Scan() {
			if val := strings.TrimSpace(scanner.Text()); val != "" {
				cfg.ScheduleEnd = val
			}
		}
	}
	
	if err := config.SaveConfig(cfg); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("\nConfiguration saved successfully!")
	showConfig()
}
