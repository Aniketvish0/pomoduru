package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/aniketvish/pomoduru/internal/config"
	"github.com/aniketvish/pomoduru/internal/timer"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	// Colors and styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1)

	timerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#04B575")).
		Padding(0, 2)

	warningStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#FF6B35")).
		Padding(0, 2)

	breakStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#00B4D8")).
		Padding(0, 2)

	extendedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#F72585")).
		Padding(0, 2)

	statusStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#4361EE")).
		Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#4CC9F0")).
		Padding(0, 2).
		Margin(0, 1)

	activeButtonStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#F72585")).
		Padding(0, 2).
		Margin(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8B8BAE")).
		Padding(1, 0)
)

// Model represents the UI model
type Model struct {
	timer       *timer.Timer
	config      *config.Config
	progress    progress.Model
	spinner     spinner.Model
	state       timer.State
	remaining   time.Duration
	width       int
	height      int
	showHelp    bool
	lastTick    time.Time
}

// NewModel creates a new UI model
func NewModel(cfg *config.Config, t *timer.Timer) Model {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F72585"))

	return Model{
		timer:     t,
		config:    cfg,
		progress:  p,
		spinner:   s,
		state:     timer.StateIdle,
		width:     maxWidth,
		height:    20,
		showHelp:  false,
		lastTick:  time.Now(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Set up timer callback
	m.timer.SetStateChangeCallback(m.handleStateChange)
	
	// Start ticking for UI updates
	return tea.Batch(
		m.spinner.Tick,
		tickCmd(),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = min(msg.Width-padding*2-4, 60)
		
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "s", " ":
			if m.state == timer.StateIdle {
				m.timer.Start()
			} else {
				m.timer.Stop()
			}
		case "e":
			if m.state == timer.StateWarning {
				if m.timer.Extend() {
					// Extension successful
				}
			}
		case "h", "?":
			m.showHelp = !m.showHelp
		}
		
	case tickMsg:
		// Update every second
		now := time.Now()
		if now.Sub(m.lastTick) >= time.Second {
			m.lastTick = now
			m.remaining = m.timer.GetRemainingTime()
			m.state = m.timer.GetState()
		}
		return m, tickCmd()
		
	case stateChangeMsg:
		m.state = timer.State(msg.state)
		m.remaining = msg.remaining
		
		// Reset progress when starting new timer
		if m.state == timer.StateWorking {
			m.progress = progress.New(
				progress.WithDefaultGradient(),
				progress.WithWidth(min(m.width-padding*2-4, 60)),
			)
		}
	}
	
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	var b strings.Builder
	
	// Title
	title := titleStyle.Render("🍅 Pomoduru - Smart Pomodoro Timer")
	b.WriteString(title + "\n\n")
	
	// Current state and time
	stateStr := m.formatState()
	timeStr := m.formatTime()
	
	var displayStr string
	switch m.state {
	case timer.StateIdle:
		displayStr = statusStyle.Render("⏸️  Ready to start")
	case timer.StateWorking:
		displayStr = timerStyle.Render("⏳ " + stateStr + " - " + timeStr)
	case timer.StateWarning:
		displayStr = warningStyle.Render("⚠️  " + stateStr + " - " + timeStr)
	case timer.StateExtended:
		displayStr = extendedStyle.Render("⏰ " + stateStr + " - " + timeStr)
	case timer.StateBreak:
		displayStr = breakStyle.Render("☕ " + stateStr + " - " + timeStr)
	case timer.StateSuspended:
		displayStr = statusStyle.Render("💤 System suspended - Taking break")
	}
	
	b.WriteString(displayStr + "\n\n")
	
	// Progress bar (only for active timers)
	if m.state == timer.StateWorking || m.state == timer.StateWarning || m.state == timer.StateExtended || m.state == timer.StateBreak {
		progress := m.calculateProgress()
		progressBar := m.progress.ViewAs(progress)
		b.WriteString(progressBar + "\n\n")
	}
	
	// Controls
	b.WriteString(m.renderControls() + "\n\n")
	
	// Info
	if m.config.AlwaysOn {
		b.WriteString(infoStyle.Render("🔄 Always-on mode: Timer will restart automatically after breaks\n"))
	}
	
	if m.config.ScheduleEnabled {
		b.WriteString(infoStyle.Render(fmt.Sprintf("📅 Scheduled: %s - %s\n", m.config.ScheduleStart, m.config.ScheduleEnd)))
	}
	
	// Help
	if m.showHelp {
		b.WriteString(m.renderHelp() + "\n")
	} else {
		b.WriteString(infoStyle.Render("Press ? for help"))
	}
	
	return lipgloss.NewStyle().
		Padding(padding).
		Width(min(m.width, maxWidth)).
		Render(b.String())
}

// Helper methods

func (m Model) formatState() string {
	switch m.state {
	case timer.StateWorking:
		return "Working"
	case timer.StateWarning:
		return "Warning"
	case timer.StateExtended:
		return "Extended"
	case timer.StateBreak:
		return "Break Time"
	default:
		return "Idle"
	}
}

func (m Model) formatTime() string {
	if m.remaining <= 0 {
		return "00:00"
	}
	
	minutes := int(m.remaining.Minutes())
	seconds := int(m.remaining.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (m Model) calculateProgress() float64 {
	var total time.Duration
	
	switch m.state {
	case timer.StateWorking, timer.StateWarning:
		total = m.config.WorkDuration
	case timer.StateExtended:
		total = m.config.ExtendDuration
	case timer.StateBreak:
		total = m.config.BreakDuration
	default:
		return 0
	}
	
	if total == 0 {
		return 0
	}
	
	elapsed := total - m.remaining
	progress := float64(elapsed) / float64(total)
	
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	
	return progress
}

func (m Model) renderControls() string {
	var controls []string
	
	if m.state == timer.StateIdle {
		controls = append(controls, activeButtonStyle.Render("[S] Start"))
	} else {
		controls = append(controls, buttonStyle.Render("[S] Stop"))
	}
	
	if m.state == timer.StateWarning && !m.timer.ExtendUsed() {
		controls = append(controls, activeButtonStyle.Render("[E] Extend (+5min)"))
	}
	
	controls = append(controls, buttonStyle.Render("[Q] Quit"))
	
	return lipgloss.JoinHorizontal(lipgloss.Top, controls...)
}

func (m Model) renderHelp() string {
	help := []string{
		"Controls:",
		"  [S] or [Space]  - Start/Stop timer",
		"  [E]             - Extend work session (5min, once per cycle)",
		"  [H] or [?]      - Toggle help",
		"  [Q] or [Ctrl+C] - Quit",
		"",
		"Features:",
		"  • Automatic system suspend after work",
		"  • 5-minute warning before suspend",
		"  • One-time 5-minute extension",
		"  • Configurable work/break durations",
		"  • Always-on mode",
		"  • Scheduled start times",
	}
	
	return infoStyle.Render(strings.Join(help, "\n"))
}

// Messages

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type stateChangeMsg struct {
	state     timer.State
	remaining time.Duration
}

func (m Model) handleStateChange(state timer.State, remaining time.Duration) {
	// This will be called from the timer, but since bubbletea is single-threaded,
	// we'll handle state changes through the tick messages
}

