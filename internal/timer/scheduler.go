package timer

import (
	"time"

	"github.com/aniketvish/pomoduru/internal/config"
)

// Scheduler manages automatic timer scheduling
type Scheduler struct {
	config *config.Config
	timer  *Timer
	active bool
}

// NewScheduler creates a new scheduler
func NewScheduler(cfg *config.Config, t *Timer) *Scheduler {
	return &Scheduler{
		config: cfg,
		timer:  t,
		active: false,
	}
}

// Start begins the scheduling system
func (s *Scheduler) Start() {
	if !s.config.ScheduleEnabled || s.active {
		return
	}
	
	s.active = true
	go s.scheduleLoop()
}

// Stop stops the scheduling system
func (s *Scheduler) Stop() {
	s.active = false
}

// scheduleLoop runs the scheduling logic
func (s *Scheduler) scheduleLoop() {
	for s.active {
		now := time.Now()
		
		// Parse schedule times
		startTime, err := time.Parse("15:04", s.config.ScheduleStart)
		if err != nil {
			// Invalid time format, skip this iteration
			time.Sleep(time.Minute)
			continue
		}
		
		endTime, err := time.Parse("15:04", s.config.ScheduleEnd)
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}
		
		// Create today's schedule times
		today := now.Truncate(24 * time.Hour)
		scheduledStart := today.Add(time.Duration(startTime.Hour())*time.Hour + time.Duration(startTime.Minute())*time.Minute)
		scheduledEnd := today.Add(time.Duration(endTime.Hour())*time.Hour + time.Duration(endTime.Minute())*time.Minute)
		
		// If end time is before start time, it means it spans midnight
		if scheduledEnd.Before(scheduledStart) {
			scheduledEnd = scheduledEnd.Add(24 * time.Hour)
		}
		
		// Check if we're within schedule
		withinSchedule := now.After(scheduledStart) && now.Before(scheduledEnd)
		
		// If we're within schedule and timer is idle, start it
		if withinSchedule && s.timer.GetState() == StateIdle && !s.config.AlwaysOn {
			s.timer.Start()
		}
		
		// If we're outside schedule and timer is running (not always-on), stop it
		if !withinSchedule && s.timer.GetState() != StateIdle && !s.config.AlwaysOn {
			// Only stop if we're not in a break or extended state
			if s.timer.GetState() == StateWorking || s.timer.GetState() == StateWarning {
				s.timer.Stop()
			}
		}
		
		// Sleep for a minute before checking again
		time.Sleep(time.Minute)
	}
}
