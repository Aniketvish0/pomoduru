# 🍅 Pomoduru - Smart Pomodoro Timer

A fancy, intelligent pomodoro timer that automatically suspends your system after work sessions, with scheduling, always-on mode, and a beautiful terminal interface.

## ✨ Features

- **Smart System Suspension**: Automatically suspends your system after work periods
- **Beautiful TUI**: Fancy terminal interface with progress bars and colors
- **Flexible Scheduling**: Set automatic start/stop times for work sessions
- **Always-On Mode**: Continuous pomodoro cycles without manual intervention
- **Extend Option**: One-time 5-minute extension when the warning appears
- **Desktop Notifications**: Get notified before system suspension
- **Systemd Integration**: Runs as a background service
- **Configurable**: Customize work/break durations, schedules, and more

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/aniketvish/pomoduru.git
cd pomoduru

# Install (requires Go)
./install.sh
```

### Basic Usage

```bash
# Run interactively
pomoduru

# Start as background service
systemctl --user start pomoduru

# Enable auto-start on boot
systemctl --user enable pomoduru
```

### Configuration

```bash
# Interactive configuration
pomoduru-config interactive

# Set specific values
pomoduru-config set --work 45 --break 15
pomoduru-config set --schedule-enabled --schedule-start 09:00 --schedule-end 18:00

# View current config
pomoduru-config show
```

## 🎮 Controls

When running interactively:

- **S** or **Space** - Start/Stop timer
- **E** - Extend work session (5min, once per cycle)
- **H** or **?** - Toggle help
- **Q** or **Ctrl+C** - Quit

## ⚙️ Configuration Options

| Setting | Default | Description |
|---------|---------|-------------|
| `work` | 50 min | Duration of work sessions |
| `break` | 10 min | Duration of break periods |
| `warning` | 5 min | Warning time before suspension |
| `extend` | 5 min | Extension duration |
| `always-on` | false | Keep timer running continuously |
| `schedule-enabled` | false | Enable scheduled start times |
| `schedule-start` | 09:00 | Automatic start time (HH:MM) |
| `schedule-end` | 18:00 | Automatic end time (HH:MM) |

## 🔧 How It Works

1. **Work Phase**: Timer counts down your work duration
2. **Warning Phase**: 5 minutes before end, shows warning and offers extension
3. **Extension**: Optional 5-minute extension (only once per cycle)
4. **Suspension**: System suspends for break period
5. **Break Phase**: After resume, break timer starts
6. **Repeat**: Cycles continue based on your settings

## 🛠️ Architecture

```
cmd/
├── pomoduru/     # Main TUI application
└── config/       # Configuration CLI tool

internal/
├── config/       # Configuration management
├── timer/        # Core timer logic + scheduler
└── ui/          # Bubbletea TUI interface

systemd/         # Systemd service files
```

## 📋 Requirements

- Linux with systemd
- Go 1.19+ (for building)
- `notify-send` (for notifications)
- `systemctl suspend` capability

## 🏗️ Building from Source

```bash
# Install dependencies
go mod tidy

# Build main application
go build -o pomoduru ./cmd/pomoduru

# Build config tool
go build -o pomoduru-config ./cmd/config
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- [Bubbletea](https://github.com/charmbracelet/bubbletea) for the beautiful TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- The Pomodoro Technique® by Francesco Cirillo

---

**Happy focusing!** 🍅
