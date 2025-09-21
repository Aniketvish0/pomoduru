package timer

import (
	"os/exec"
	"time"

	"github.com/aniketvish/pomoduru/internal/config"
)

// State represents the current timer state
type State int

const (
	StateIdle State = iota
	StateWorking
	StateWarning
	StateBreak
	StateExtended
	StateSuspended
)

// Timer manages the pomodoro timer
type Timer struct {
	config        *config.Config
	state         State
	startTime     time.Time
	extendUsed    bool // Track if extension has been used this cycle
	onStateChange func(State, time.Duration) // Callback for state changes
}

// NewTimer creates a new timer instance
func NewTimer(cfg *config.Config) *Timer {
	return &Timer{
		config:     cfg,
		state:      StateIdle,
		extendUsed: false,
	}
}

// SetStateChangeCallback sets the callback for state changes
func (t *Timer) SetStateChangeCallback(callback func(State, time.Duration)) {
	t.onStateChange = callback
}

// Start begins the pomodoro timer
func (t *Timer) Start() {
	t.state = StateWorking
	t.startTime = time.Now()
	t.extendUsed = false
	
	if t.onStateChange != nil {
		t.onStateChange(t.state, t.config.WorkDuration)
	}
	
	// Start the warning timer
	warningDuration := t.config.WorkDuration - t.config.WarningTime
	time.AfterFunc(warningDuration, t.handleWarning)
	
	// Start the work timer
	time.AfterFunc(t.config.WorkDuration, t.handleWorkComplete)
}

// Extend extends the current work session by the configured extend duration
// Returns true if extension was allowed, false if already used
func (t *Timer) Extend() bool {
	if t.extendUsed || t.state != StateWarning {
		return false
	}
	
	t.extendUsed = true
	t.state = StateExtended
	
	// Cancel the existing timers and create new ones
	time.AfterFunc(t.config.ExtendDuration, t.handleExtendedWorkComplete)
	
	if t.onStateChange != nil {
		t.onStateChange(t.state, t.config.ExtendDuration)
	}
	
	return true
}

// Stop stops the timer and resets to idle state
func (t *Timer) Stop() {
	t.state = StateIdle
	t.extendUsed = false
	
	if t.onStateChange != nil {
		t.onStateChange(t.state, 0)
	}
}

// GetState returns the current timer state
func (t *Timer) GetState() State {
	return t.state
}

// GetRemainingTime returns the remaining time in current state
func (t *Timer) GetRemainingTime() time.Duration {
	if t.state == StateIdle {
		return 0
	}
	
	var totalDuration time.Duration
	switch t.state {
	case StateWorking, StateWarning:
		totalDuration = t.config.WorkDuration
	case StateExtended:
		totalDuration = t.config.ExtendDuration
	case StateBreak:
		totalDuration = t.config.BreakDuration
	}
	
	elapsed := time.Since(t.startTime)
	return totalDuration - elapsed
}

// ExtendUsed returns whether extension has been used this cycle
func (t *Timer) ExtendUsed() bool {
	return t.extendUsed
}

// handleWarning is called when warning time is reached
func (t *Timer) handleWarning() {
	if t.state == StateWorking {
		t.state = StateWarning
		
		// Send notification
		exec.Command("notify-send", "Pomoduru", "System will sleep in 5 minutes! Use 'Extend' to delay.").Run()
		
		if t.onStateChange != nil {
			t.onStateChange(t.state, t.config.WarningTime)
		}
	}
}

// handleWorkComplete is called when work time is complete
func (t *Timer) handleWorkComplete() {
	if t.state == StateWorking || t.state == StateWarning {
		t.suspendSystem()
	}
}

// handleExtendedWorkComplete is called when extended work time is complete
func (t *Timer) handleExtendedWorkComplete() {
	if t.state == StateExtended {
		t.suspendSystem()
	}
}

// suspendSystem suspends the system
func (t *Timer) suspendSystem() {
	t.state = StateSuspended
	
	if t.onStateChange != nil {
		t.onStateChange(t.state, 0)
	}
	
	// Send final notification
	exec.Command("notify-send", "Pomoduru", "Time's up! Taking a break...").Run()
	
	// Suspend system
	exec.Command("systemctl", "suspend", "-i").Run()
	
	// After suspend and resume, start break timer
	time.AfterFunc(time.Second, func() {
		t.startBreak()
	})
}

// startBreak starts the break period
func (t *Timer) startBreak() {
	t.state = StateBreak
	t.startTime = time.Now()
	
	if t.onStateChange != nil {
		t.onStateChange(t.state, t.config.BreakDuration)
	}
	
	// Start break timer
	time.AfterFunc(t.config.BreakDuration, t.handleBreakComplete)
}

// handleBreakComplete is called when break time is complete
func (t *Timer) handleBreakComplete() {
	if t.state == StateBreak {
		// If always-on mode, restart the cycle
		if t.config.AlwaysOn {
			t.Start()
		} else {
			t.Stop()
		}
	}
}
