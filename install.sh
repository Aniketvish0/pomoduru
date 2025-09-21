#!/bin/bash

set -e

echo "🍅 Installing Pomoduru..."

# Build the binaries
echo "Building binaries..."
go build -o pomoduru ./cmd/pomoduru
go build -o pomoduru-config ./cmd/config

# Install binaries
echo "Installing binaries to /usr/local/bin..."
sudo cp pomoduru /usr/local/bin/
sudo cp pomoduru-config /usr/local/bin/
sudo chmod +x /usr/local/bin/pomoduru
sudo chmod +x /usr/local/bin/pomoduru-config

# Install systemd service
echo "Installing systemd service..."
mkdir -p ~/.config/systemd/user
cp systemd/pomoduru.service ~/.config/systemd/user/

# Reload systemd
systemctl --user daemon-reload

echo "✅ Installation complete!"
echo ""
echo "To start Pomoduru:"
echo "  systemctl --user start pomoduru"
echo ""
echo "To start Pomoduru on boot:"
echo "  systemctl --user enable pomoduru"
echo ""
echo "To configure Pomoduru:"
echo "  pomoduru-config interactive"
echo ""
echo "To run interactively:"
echo "  pomoduru"
echo ""
echo "Happy focusing! 🍅"
